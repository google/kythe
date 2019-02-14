/*
 * Copyright 2014 The Kythe Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Defines AST visitors and consumers used by the indexer tool.

#ifndef KYTHE_CXX_INDEXER_CXX_INDEXER_AST_HOOKS_H_
#define KYTHE_CXX_INDEXER_CXX_INDEXER_AST_HOOKS_H_

#include <memory>
#include <unordered_map>
#include <unordered_set>
#include <utility>

#include "absl/memory/memory.h"
#include "absl/types/optional.h"
#include "clang/AST/ASTContext.h"
#include "clang/AST/ASTTypeTraits.h"
#include "clang/AST/RecursiveASTVisitor.h"
#include "clang/Index/USRGeneration.h"
#include "clang/Sema/SemaConsumer.h"
#include "clang/Sema/Template.h"
#include "glog/logging.h"

#include "GraphObserver.h"
#include "IndexerLibrarySupport.h"
#include "indexed_parent_map.h"
#include "indexer_worklist.h"
#include "kythe/cxx/indexer/cxx/node_set.h"
#include "kythe/cxx/indexer/cxx/semantic_hash.h"
#include "marked_source.h"
#include "type_map.h"

namespace kythe {

/// \brief Specifies whether uncommonly-used data should be dropped.
enum Verbosity : bool {
  Classic = true,  ///< Emit all data.
  Lite = false     ///< Emit only common data.
};

/// \brief Specifies what the indexer should do if it encounters a case it
/// doesn't understand.
enum BehaviorOnUnimplemented : bool {
  Abort = false,   ///< Stop indexing and exit with an error.
  Continue = true  ///< Continue indexing, possibly emitting less data.
};

/// \brief Specifies what the indexer should do with template instantiations.
enum BehaviorOnTemplates : bool {
  SkipInstantiations = false,  ///< Don't visit template instantiations.
  VisitInstantiations = true   ///< Visit template instantiations.
};

/// \brief Specifies if the indexer should emit documentation nodes for comments
/// associated with forward declarations.
enum BehaviorOnFwdDeclComments : bool { Emit = true, Ignore = false };

/// \brief A byte range that links to some node.
struct MiniAnchor {
  size_t Begin;
  size_t End;
  GraphObserver::NodeId AnchoredTo;
};

/// Adds brackets to Text to define anchor locations (escaping existing ones)
/// and sorts Anchors such that the ith Anchor corresponds to the ith opening
/// bracket. Drops empty or negative-length spans.
void InsertAnchorMarks(std::string& Text, std::vector<MiniAnchor>& Anchors);

/// \brief Used internally to check whether parts of the AST can be ignored.
class PruneCheck;

/// \brief An AST visitor that extracts information for a translation unit and
/// writes it to a `GraphObserver`.
class IndexerASTVisitor : public clang::RecursiveASTVisitor<IndexerASTVisitor> {
 public:
  IndexerASTVisitor(clang::ASTContext& C, BehaviorOnUnimplemented B,
                    BehaviorOnTemplates T, Verbosity V,
                    BehaviorOnFwdDeclComments ObjC,
                    BehaviorOnFwdDeclComments Cpp, const LibrarySupports& S,
                    clang::Sema& Sema, std::function<bool()> ShouldStopIndexing,
                    GraphObserver* GO = nullptr, int UsrByteSize = 0)
      : IgnoreUnimplemented(B),
        TemplateMode(T),
        Verbosity(V),
        ObjCFwdDocs(ObjC),
        CppFwdDocs(Cpp),
        Observer(GO ? *GO : NullObserver),
        Context(C),
        Supports(S),
        Sema(Sema),
        MarkedSources(&Sema, &Observer),
        ShouldStopIndexing(std::move(ShouldStopIndexing)),
        UsrByteSize(UsrByteSize) {}

  bool VisitDecl(const clang::Decl* Decl);
  bool VisitFieldDecl(const clang::FieldDecl* Decl);
  bool VisitVarDecl(const clang::VarDecl* Decl);
  bool VisitNamespaceDecl(const clang::NamespaceDecl* Decl);
  bool VisitBindingDecl(const clang::BindingDecl* Decl);
  bool VisitDeclRefExpr(const clang::DeclRefExpr* DRE);
  bool VisitDesignatedInitExpr(const clang::DesignatedInitExpr* DIE);
  bool VisitCXXConstructExpr(const clang::CXXConstructExpr* E);
  bool VisitCXXDeleteExpr(const clang::CXXDeleteExpr* E);
  bool VisitCXXNewExpr(const clang::CXXNewExpr* E);
  bool VisitCXXPseudoDestructorExpr(const clang::CXXPseudoDestructorExpr* E);
  bool VisitCXXUnresolvedConstructExpr(
      const clang::CXXUnresolvedConstructExpr* E);
  bool VisitCallExpr(const clang::CallExpr* Expr);
  bool VisitMemberExpr(const clang::MemberExpr* Expr);
  bool VisitCXXDependentScopeMemberExpr(
      const clang::CXXDependentScopeMemberExpr* Expr);

  // Visitors for leaf TypeLoc types.
  bool VisitBuiltinTypeLoc(clang::BuiltinTypeLoc TL);
  bool VisitEnumTypeLoc(clang::EnumTypeLoc TL);
  bool VisitRecordTypeLoc(clang::RecordTypeLoc TL);
  bool VisitTemplateTypeParmTypeLoc(clang::TemplateTypeParmTypeLoc TL);
  bool VisitSubstTemplateTypeParmTypeLoc(
      clang::SubstTemplateTypeParmTypeLoc TL);
  bool VisitTemplateSpecializationTypeLoc(
      clang::TemplateSpecializationTypeLoc TL);
  // Handles AutoTypeLoc and DeducedTemplateSpecializationTypeLoc
  bool VisitDeducedTypeLoc(clang::DeducedTypeLoc TL);
  bool VisitDecltypeTypeLoc(clang::DecltypeTypeLoc TL);
  bool VisitElaboratedTypeLoc(clang::ElaboratedTypeLoc TL);
  bool VisitTypedefTypeLoc(clang::TypedefTypeLoc TL);
  bool VisitInjectedClassNameTypeLoc(clang::InjectedClassNameTypeLoc TL);
  bool VisitDependentNameTypeLoc(clang::DependentNameTypeLoc TL);
  bool VisitPackExpansionTypeLoc(clang::PackExpansionTypeLoc TL);
  bool VisitObjCObjectTypeLoc(clang::ObjCObjectTypeLoc TL);
  bool VisitObjCTypeParamTypeLoc(clang::ObjCTypeParamTypeLoc TL);

  bool TraverseAttributedTypeLoc(clang::AttributedTypeLoc TL);
  bool TraverseDependentAddressSpaceTypeLoc(
      clang::DependentAddressSpaceTypeLoc TL);

  bool TraverseMemberPointerTypeLoc(clang::MemberPointerTypeLoc TL);

  // Emit edges for an anchor pointing to the indicated type.
  NodeSet RecordTypeLocSpellingLocation(clang::TypeLoc TL);

  bool TraverseDeclarationNameInfo(clang::DeclarationNameInfo NameInfo);

  // Visit the subtypes of TypedefNameDecl individually because we want to do
  // something different with ObjCTypeParamDecl.
  bool VisitTypedefDecl(const clang::TypedefDecl* Decl);
  bool VisitTypeAliasDecl(const clang::TypeAliasDecl* Decl);
  bool VisitObjCTypeParamDecl(const clang::ObjCTypeParamDecl* Decl);
  bool VisitUsingShadowDecl(const clang::UsingShadowDecl* Decl);

  bool VisitRecordDecl(const clang::RecordDecl* Decl);
  bool VisitEnumDecl(const clang::EnumDecl* Decl);
  bool VisitEnumConstantDecl(const clang::EnumConstantDecl* Decl);
  bool VisitFunctionDecl(clang::FunctionDecl* Decl);
  bool TraverseDecl(clang::Decl* Decl);

  bool TraverseConstructorInitializer(clang::CXXCtorInitializer* Init);
  bool TraverseCXXNewExpr(clang::CXXNewExpr* E);
  bool TraverseCXXFunctionalCastExpr(clang::CXXFunctionalCastExpr* E);

  bool IndexConstructExpr(const clang::CXXConstructExpr* E,
                          const clang::TypeSourceInfo* TSI);

  // Objective C specific nodes
  bool VisitObjCPropertyImplDecl(const clang::ObjCPropertyImplDecl* Decl);
  bool VisitObjCCompatibleAliasDecl(const clang::ObjCCompatibleAliasDecl* Decl);
  bool VisitObjCCategoryDecl(const clang::ObjCCategoryDecl* Decl);
  bool VisitObjCImplementationDecl(
      const clang::ObjCImplementationDecl* ImplDecl);
  bool VisitObjCCategoryImplDecl(const clang::ObjCCategoryImplDecl* ImplDecl);
  bool VisitObjCInterfaceDecl(const clang::ObjCInterfaceDecl* Decl);
  bool VisitObjCProtocolDecl(const clang::ObjCProtocolDecl* Decl);
  bool VisitObjCMethodDecl(const clang::ObjCMethodDecl* Decl);
  bool VisitObjCPropertyDecl(const clang::ObjCPropertyDecl* Decl);
  bool VisitObjCIvarRefExpr(const clang::ObjCIvarRefExpr* IRE);
  bool VisitObjCMessageExpr(const clang::ObjCMessageExpr* Expr);
  bool VisitObjCPropertyRefExpr(const clang::ObjCPropertyRefExpr* Expr);

  // TODO(salguarnieri) We could link something here (the square brackets?) to
  // the setter and getter methods
  //   bool VisitObjCSubscriptRefExpr(const clang::ObjCSubscriptRefExpr *Expr);

  // TODO(salguarnieri) delete this comment block when we have more objective-c
  // support implemented.
  //
  //  Visitors that are left to their default behavior because we do not need
  //  to take any action while visiting them.
  //   bool VisitObjCDictionaryLiteral(const clang::ObjCDictionaryLiteral *D);
  //   bool VisitObjCArrayLiteral(const clang::ObjCArrayLiteral *D);
  //   bool VisitObjCBoolLiteralExpr(const clang::ObjCBoolLiteralExpr *D);
  //   bool VisitObjCStringLiteral(const clang::ObjCStringLiteral *D);
  //   bool VisitObjCEncodeExpr(const clang::ObjCEncodeExpr *Expr);
  //   bool VisitObjCBoxedExpr(const clang::ObjCBoxedExpr *Expr);
  //   bool VisitObjCSelectorExpr(const clang::ObjCSelectorExpr *Expr);
  //   bool VisitObjCIndirectCopyRestoreExpr(
  //    const clang::ObjCIndirectCopyRestoreExpr *Expr);
  //   bool VisitObjCIsaExpr(const clang::ObjCIsaExpr *Expr);
  //
  //  We visit the subclasses of ObjCContainerDecl so there is nothing to do.
  //   bool VisitObjCContainerDecl(const clang::ObjCContainerDecl *D);
  //
  //  We visit the subclasses of ObjCImpleDecl so there is nothing to do.
  //   bool VisitObjCImplDecl(const clang::ObjCImplDecl *D);
  //
  //  There is nothing specific we need to do for blocks. The recursive visitor
  //  will visit the components of them correctly.
  //   bool VisitBlockDecl(const clang::BlockDecl *Decl);
  //   bool VisitBlockExpr(const clang::BlockExpr *Expr);

  /// \brief For functions that support it, controls the emission of range
  /// information.
  enum class EmitRanges {
    No,  ///< Don't emit range information.
    Yes  ///< Emit range information when it's available.
  };

  // Objective C methods don't have TypeSourceInfo so we must construct a type
  // for the methods to be used in the graph.
  absl::optional<GraphObserver::NodeId> CreateObjCMethodTypeNode(
      const clang::ObjCMethodDecl* MD, EmitRanges ER);

  /// \brief Builds a stable node ID for a compile-time expression.
  /// \param Expr The expression to represent.
  /// \param ER Whether to notify the `GraphObserver` about source text
  /// ranges for expressions.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForExpr(
      const clang::Expr* Expr, EmitRanges ER);

  /// \brief Builds a stable node ID for a special template argument.
  /// \param Id A string representing the special argument.
  GraphObserver::NodeId BuildNodeIdForSpecialTemplateArgument(
      llvm::StringRef Id);

  /// \brief Builds a stable node ID for a template expansion template argument.
  /// \param Name The template pattern being expanded.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForTemplateExpansion(
      clang::TemplateName Name);

  /// \brief Builds a stable NodeSet for the given TypeLoc.
  /// \param TL The TypeLoc for which to build a NodeSet.
  /// \returns NodeSet instance indicating claimability of the contained
  /// NodeIds.
  NodeSet BuildNodeSetForType(const clang::TypeLoc& TL);
  NodeSet BuildNodeSetForType(const clang::QualType& QT);

  NodeSet BuildNodeSetForBuiltin(clang::BuiltinTypeLoc TL) const;
  NodeSet BuildNodeSetForEnum(clang::EnumTypeLoc TL);
  NodeSet BuildNodeSetForRecord(clang::RecordTypeLoc TL);
  NodeSet BuildNodeSetForInjectedClassName(clang::InjectedClassNameTypeLoc TL);
  NodeSet BuildNodeSetForTemplateTypeParm(clang::TemplateTypeParmTypeLoc TL);
  NodeSet BuildNodeSetForPointer(clang::PointerTypeLoc TL);
  NodeSet BuildNodeSetForMemberPointer(clang::MemberPointerTypeLoc TL);
  NodeSet BuildNodeSetForLValueReference(clang::LValueReferenceTypeLoc TL);
  NodeSet BuildNodeSetForRValueReference(clang::RValueReferenceTypeLoc TL);

  NodeSet BuildNodeSetForAuto(clang::AutoTypeLoc TL);
  NodeSet BuildNodeSetForDeducedTemplateSpecialization(
      clang::DeducedTemplateSpecializationTypeLoc TL);

  NodeSet BuildNodeSetForQualified(clang::QualifiedTypeLoc TL);
  NodeSet BuildNodeSetForConstantArray(clang::ConstantArrayTypeLoc TL);
  NodeSet BuildNodeSetForIncompleteArray(clang::IncompleteArrayTypeLoc TL);
  NodeSet BuildNodeSetForDependentSizedArray(
      clang::DependentSizedArrayTypeLoc TL);
  NodeSet BuildNodeSetForFunctionProto(clang::FunctionProtoTypeLoc TL);
  NodeSet BuildNodeSetForFunctionNoProto(clang::FunctionNoProtoTypeLoc TL);
  NodeSet BuildNodeSetForParen(clang::ParenTypeLoc TL);
  NodeSet BuildNodeSetForDecltype(clang::DecltypeTypeLoc TL);
  NodeSet BuildNodeSetForElaborated(clang::ElaboratedTypeLoc TL);
  NodeSet BuildNodeSetForTypedef(clang::TypedefTypeLoc TL);

  NodeSet BuildNodeSetForSubstTemplateTypeParm(
      clang::SubstTemplateTypeParmTypeLoc TL);
  NodeSet BuildNodeSetForDependentName(clang::DependentNameTypeLoc TL);
  NodeSet BuildNodeSetForTemplateSpecialization(
      clang::TemplateSpecializationTypeLoc TL);
  NodeSet BuildNodeSetForPackExpansion(clang::PackExpansionTypeLoc TL);
  NodeSet BuildNodeSetForBlockPointer(clang::BlockPointerTypeLoc TL);
  NodeSet BuildNodeSetForObjCObjectPointer(clang::ObjCObjectPointerTypeLoc TL);
  NodeSet BuildNodeSetForObjCObject(clang::ObjCObjectTypeLoc TL);
  NodeSet BuildNodeSetForObjCTypeParam(clang::ObjCTypeParamTypeLoc TL);
  NodeSet BuildNodeSetForObjCInterface(clang::ObjCInterfaceTypeLoc TL);
  NodeSet BuildNodeSetForAttributed(clang::AttributedTypeLoc TL);
  NodeSet BuildNodeSetForDependentAddressSpace(
      clang::DependentAddressSpaceTypeLoc TL);

  // Helper used for Auto and DeducedTemplateSpecialization.
  NodeSet BuildNodeSetForDeduced(clang::DeducedTypeLoc TL);

  // Helper function which constructs marked source and records
  // a tnominal node for the given `Decl`.
  GraphObserver::NodeId BuildNominalNodeIdForDecl(const clang::NamedDecl* Decl);

  // Helper used by BuildNodeSetForRecord and BuildNodeSetForInjectedClassName.
  NodeSet BuildNodeSetForNonSpecializedRecordDecl(
      const clang::RecordDecl* Decl);

  const clang::TemplateTypeParmDecl* FindTemplateTypeParmTypeLocDecl(
      clang::TemplateTypeParmTypeLoc TL) const;

  absl::optional<GraphObserver::NodeId> BuildNodeIdForObjCProtocols(
      clang::ObjCObjectTypeLoc TL);
  GraphObserver::NodeId BuildNodeIdForObjCProtocols(
      const clang::ObjCObjectType* T);
  GraphObserver::NodeId BuildNodeIdForObjCProtocols(
      GraphObserver::NodeId BaseType, const clang::ObjCObjectType* T);

  /// \brief Builds a stable node ID for `Type`.
  /// \param TypeLoc The type that is being identified.
  /// \return The Node ID for `Type`.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForType(
      const clang::TypeLoc& TypeLoc);

  /// \brief Builds a stable node ID for `QT`.
  /// \param QT The type that is being identified.
  /// \return The Node ID for `QT`.
  ///
  /// This function will invent a `TypeLoc` with an invalid location.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForType(
      const clang::QualType& QT);

  /// \brief Builds a stable node ID for the given `TemplateName`.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForTemplateName(
      const clang::TemplateName& Name);

  /// \brief Builds a stable node ID for the given `TemplateArgument`.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForTemplateArgument(
      const clang::TemplateArgumentLoc& Arg, EmitRanges ER);

  /// \brief Builds a stable node ID for the given `TemplateArgument`.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForTemplateArgument(
      const clang::TemplateArgument& Arg, clang::SourceLocation L);

  /// \brief Builds a stable node ID for `Stmt`.
  ///
  /// This mechanism for generating IDs should only be used in contexts where
  /// identifying statements through source locations/wraith contexts is not
  /// possible (e.g., in implicit code).
  ///
  /// \param Decl The statement that is being identified
  /// \return The node for `Stmt` if the statement was implicit; otherwise,
  /// None.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForImplicitStmt(
      const clang::Stmt* Stmt);

  /// \brief Builds a stable node ID for `Decl`'s tapp if it's an implicit
  /// template instantiation.
  absl::optional<GraphObserver::NodeId>
  BuildNodeIdForImplicitTemplateInstantiation(const clang::Decl* Decl);

  /// \brief Builds a stable node ID for `Decl`'s tapp if it's an implicit
  /// function template instantiation.
  absl::optional<GraphObserver::NodeId>
  BuildNodeIdForImplicitFunctionTemplateInstantiation(
      const clang::FunctionDecl* Decl);

  /// \brief Builds a stable node ID for `Decl`'s tapp if it's an implicit
  /// class template instantiation.
  absl::optional<GraphObserver::NodeId>
  BuildNodeIdForImplicitClassTemplateInstantiation(
      const clang::ClassTemplateSpecializationDecl* Decl);

  /// \brief Builds a stable node ID for an external reference to `Decl`.
  ///
  /// This is equivalent to BuildNodeIdForDecl for Decls that are not
  /// implicit template instantiations; otherwise, it returns the `NodeId`
  /// for the tapp node for the instantiation.
  ///
  /// \param Decl The declaration that is being identified.
  /// \return The node for `Decl`.
  GraphObserver::NodeId BuildNodeIdForRefToDecl(const clang::Decl* Decl);

  /// \brief Builds a stable node ID for `Decl`.
  ///
  /// \param Decl The declaration that is being identified.
  /// \return The node for `Decl`.
  GraphObserver::NodeId BuildNodeIdForDecl(const clang::Decl* Decl);

  /// \brief Builds a stable node ID for the definition of `Decl`, if
  /// `Decl` is not already a definition and its definition can be found.
  ///
  /// \param Decl The declaration that is being identified.
  /// \return The node for the definition `Decl` if `Decl` isn't a definition
  /// and its definition can be found; or None.
  template <typename D>
  absl::optional<GraphObserver::NodeId> BuildNodeIdForDefnOfDecl(
      const D* Decl) {
    if (const auto* Defn = Decl->getDefinition()) {
      if (Defn != Decl) {
        return BuildNodeIdForDecl(Defn);
      }
    }
    return absl::nullopt;
  }

  /// \brief Builds a stable node ID for `TND`.
  ///
  /// \param Decl The declaration that is being identified.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForTypedefNameDecl(
      const clang::TypedefNameDecl* TND);

  /// \brief Builds a stable node ID for `Decl`.
  ///
  /// There is not a one-to-one correspondence between `Decl`s and nodes.
  /// Certain `Decl`s, like `TemplateTemplateParmDecl`, are split into a
  /// node representing the parameter and a node representing the kind of
  /// the abstraction. The primary node is returned by the normal
  /// `BuildNodeIdForDecl` function.
  ///
  /// \param Decl The declaration that is being identified.
  /// \param Index The index of the sub-id to generate.
  ///
  /// \return A stable node ID for `Decl`'s `Index`th subnode.
  GraphObserver::NodeId BuildNodeIdForDecl(const clang::Decl* Decl,
                                           unsigned Index);

  /// \brief Categorizes the name of `Decl` according to the equivalence classes
  /// defined by `GraphObserver::NameId::NameEqClass`.
  GraphObserver::NameId::NameEqClass BuildNameEqClassForDecl(
      const clang::Decl* Decl) const;

  /// \brief Builds a stable name ID for the name of `Decl`.
  ///
  /// \param Decl The declaration that is being named.
  /// \return The name for `Decl`.
  GraphObserver::NameId BuildNameIdForDecl(const clang::Decl* Decl);

  /// \brief Builds a NodeId for the given dependent name.
  ///
  /// \param NNS The qualifier on the name.
  /// \param Id The name itself.
  /// \param IdLoc The name's location.
  /// \param Root If provided, the primary NodeId is morally prepended to `NNS`
  /// such that the dependent name is lookup(lookup*(Root, NNS), Id).
  absl::optional<GraphObserver::NodeId> BuildNodeIdForDependentName(
      const clang::NestedNameSpecifierLoc& NNS,
      const clang::DeclarationName& Id, const clang::SourceLocation IdLoc,
      const absl::optional<GraphObserver::NodeId>& Root);

  GraphObserver::NodeId BuildNodeIdForDependentLoc(
      const clang::NestedNameSpecifierLoc& NNSLoc,
      const clang::SourceLocation& IdLoc);

  GraphObserver::NodeId BuildNodeIdForDependentRange(
      const clang::NestedNameSpecifierLoc& NNSLoc,
      const clang::SourceRange& IdRange);

  absl::optional<GraphObserver::NodeId> RecordDependentParamEdges();
  GraphObserver::NodeId RecordDependentLookup(
      const GraphObserver::NodeId& DID, const clang::DeclarationName& Name);

  bool TraverseNestedNameSpecifierLoc(clang::NestedNameSpecifierLoc NNS);

  /// \brief Is `VarDecl` a definition?
  ///
  /// A parameter declaration is considered a definition if it appears as part
  /// of a function definition; otherwise it is always a declaration. This
  /// differs from the C++ Standard, which treats these parameters as
  /// definitions (basic.scope.proto).
  static bool IsDefinition(const clang::VarDecl* VD);

  /// \brief Is `FunctionDecl` a definition?
  static bool IsDefinition(const clang::FunctionDecl* FunctionDecl);

  /// \brief Gets a range for the name of a declaration, whether that name is a
  /// single token or otherwise.
  ///
  /// The returned range is a best-effort attempt to cover the "name" of
  /// the entity as written in the source code.
  clang::SourceRange RangeForNameOfDeclaration(
      const clang::NamedDecl* Decl) const;

  /// \brief Gets a suitable range for an AST entity from the `start_location`.
  clang::SourceRange RangeForASTEntity(
      clang::SourceLocation start_location) const;
  clang::SourceRange RangeForSingleToken(
      clang::SourceLocation start_location) const;

  /// Consume a token of the `ExpectedKind` from the `StartLocation`,
  /// returning the range for that token on success and an invalid
  /// range otherwise.
  ///
  /// The begin location for the returned range may be different than
  /// StartLocation. For example, this can happen if StartLocation points to
  /// whitespace before the start of the token.
  clang::SourceRange ConsumeToken(clang::SourceLocation StartLocation,
                                  clang::tok::TokenKind ExpectedKind) const;

  bool TraverseClassTemplateDecl(clang::ClassTemplateDecl* TD);
  bool TraverseClassTemplateSpecializationDecl(
      clang::ClassTemplateSpecializationDecl* TD);
  bool TraverseClassTemplatePartialSpecializationDecl(
      clang::ClassTemplatePartialSpecializationDecl* TD);

  bool TraverseVarTemplateDecl(clang::VarTemplateDecl* TD);
  bool TraverseVarTemplateSpecializationDecl(
      clang::VarTemplateSpecializationDecl* VD);
  bool ForceTraverseVarTemplateSpecializationDecl(
      clang::VarTemplateSpecializationDecl* VD);
  bool TraverseVarTemplatePartialSpecializationDecl(
      clang::VarTemplatePartialSpecializationDecl* TD);

  bool TraverseFunctionDecl(clang::FunctionDecl* FD);
  bool TraverseFunctionTemplateDecl(clang::FunctionTemplateDecl* FTD);

  bool TraverseTypeAliasTemplateDecl(clang::TypeAliasTemplateDecl* TATD);

  bool shouldVisitTemplateInstantiations() const {
    return TemplateMode == BehaviorOnTemplates::VisitInstantiations;
  }
  bool shouldEmitObjCForwardClassDeclDocumentation() const {
    return ObjCFwdDocs == BehaviorOnFwdDeclComments::Emit;
  }
  bool shouldEmitCppForwardDeclDocumentation() const {
    return CppFwdDocs == BehaviorOnFwdDeclComments::Emit;
  }
  bool shouldVisitImplicitCode() const { return true; }
  // Disables data recursion. We intercept Traverse* methods in the RAV, which
  // are not triggered during data recursion.
  bool shouldUseDataRecursionFor(clang::Stmt* S) const { return false; }

  /// \brief Records the range of a definition. If the range cannot be placed
  /// somewhere inside a source file, no record is made.
  void MaybeRecordDefinitionRange(
      const absl::optional<GraphObserver::Range>& R,
      const GraphObserver::NodeId& Id,
      const absl::optional<GraphObserver::NodeId>& DeclId);

  /// \brief Returns the attached GraphObserver.
  GraphObserver& getGraphObserver() { return Observer; }

  /// \brief Returns the ASTContext.
  const clang::ASTContext& getASTContext() { return Context; }

  /// If `SR` is empty (getBegin() == getEnd()) and a valid file id, expands the
  /// range. Otherwise, returns the input unmodified.
  clang::SourceRange ExpandRangeIfEmptyFileID(const clang::SourceRange& SR);

  // If `SR` is a valid macro id, attempt to map it to a file range,
  // otherwise returns the input unmodified.
  clang::SourceRange MapRangeToFileIfMacroID(const clang::SourceRange& SR);

  /// Returns `SR` as a `Range` in this `RecursiveASTVisitor`'s current
  /// RangeContext after expanding empty ranges and mapping macros to a file
  /// location.
  absl::optional<GraphObserver::Range> ExpandedFileRangeInCurrentContext(
      const clang::SourceRange& SR);

  /// Returns `SR` as a `Range` in this `RecursiveASTVisitor`'s current
  /// RangeContext.
  absl::optional<GraphObserver::Range> ExplicitRangeInCurrentContext(
      const clang::SourceRange& SR);

  /// Returns `SR` as a `Range` in this `RecursiveASTVisitor`'s current
  /// RangeContext. If SR is in a macro, the returned Range will be mapped
  /// to a file first. If the range would be zero-width, it will be expanded
  /// via RangeForASTEntityFromSourceLocation.
  absl::optional<GraphObserver::Range> ExpandedRangeInCurrentContext(
      clang::SourceRange SR);

  /// If `Implicit` is true, returns `Id` as an implicit Range; otherwise,
  /// returns `SR` as a `Range` in this `RecursiveASTVisitor`'s current
  /// RangeContext.
  absl::optional<GraphObserver::Range> RangeInCurrentContext(
      bool Implicit, const GraphObserver::NodeId& Id,
      const clang::SourceRange& SR);

  /// If `Id` is some NodeId, returns it as an implicit Range; otherwise,
  /// returns `SR` as a `Range` in this `RecursiveASTVisitor`'s current
  /// RangeContext.
  absl::optional<GraphObserver::Range> RangeInCurrentContext(
      const absl::optional<GraphObserver::NodeId>& Id,
      const clang::SourceRange& SR);

  void RunJob(std::unique_ptr<IndexJob> JobToRun) {
    Job = std::move(JobToRun);
    if (Job->IsDeclJob) {
      TraverseDecl(Job->Decl);
    } else {
      // There is no declaration attached to a top-level file comment.
      HandleFileLevelComments(Job->FileId, Job->FileNode);
    }
  }

  const IndexJob& getCurrentJob() {
    CHECK(Job != nullptr);
    return *Job;
  }

  void Work(clang::Decl* InitialDecl,
            std::unique_ptr<IndexerWorklist> NewWorklist) {
    Worklist = std::move(NewWorklist);
    Worklist->EnqueueJob(llvm::make_unique<IndexJob>(InitialDecl));
    while (!ShouldStopIndexing() && Worklist->DoWork())
      ;
    Observer.iterateOverClaimedFiles(
        [this, InitialDecl](clang::FileID Id,
                            const GraphObserver::NodeId& FileNode) {
          RunJob(llvm::make_unique<IndexJob>(InitialDecl, Id, FileNode));
          return !ShouldStopIndexing();
        });
    Worklist.reset();
  }

  /// \brief Provides execute-only access to ShouldStopIndexing. Should be used
  /// from the same thread that's walking the AST.
  bool shouldStopIndexing() const { return ShouldStopIndexing(); }

  /// Blames a call to `Callee` at `Range` on everything at the top of
  /// `BlameStack` (or does nothing if there's nobody to blame).
  void RecordCallEdges(const GraphObserver::Range& Range,
                       const GraphObserver::NodeId& Callee);

  /// \return whether `range` should be considered to be implicit under the
  /// current context.
  GraphObserver::Implicit IsImplicit(const GraphObserver::Range& range);

 private:
  using Base = RecursiveASTVisitor;

  friend class PruneCheck;

  /// Whether we should stop on missing cases or continue on.
  BehaviorOnUnimplemented IgnoreUnimplemented;

  /// Should we visit template instantiations?
  BehaviorOnTemplates TemplateMode;

  /// Should we emit all data?
  enum Verbosity Verbosity;

  /// Should we emit documentation for forward class decls in ObjC?
  BehaviorOnFwdDeclComments ObjCFwdDocs;

  /// Should we emit documentation for forward decls in C++?
  BehaviorOnFwdDeclComments CppFwdDocs;

  NullGraphObserver NullObserver;
  GraphObserver& Observer;
  clang::ASTContext& Context;

  /// \brief The result of calling into the lexer.
  enum class LexerResult {
    Failure,  ///< The operation failed.
    Success   ///< The operation completed.
  };

  /// \brief Using the `Observer`'s preprocessor, relexes the token at the
  /// specified location. Ignores whitespace.
  /// \param StartLocation Where to begin lexing.
  /// \param Token The token to overwrite.
  /// \return `Failure` if there was a failure, `Success` on success.
  LexerResult getRawToken(clang::SourceLocation StartLocation,
                          clang::Token& Token) const {
    return Observer.getPreprocessor()->getRawToken(StartLocation, Token,
                                                   true /* ignoreWhiteSpace */)
               ? LexerResult::Failure
               : LexerResult::Success;
  }

  /// \brief Handle the file-level comments for `Id` with node ID `FileId`.
  void HandleFileLevelComments(clang::FileID Id,
                               const GraphObserver::NodeId& FileId);

  /// \brief Emit data for `Comment` that documents `DocumentedNode`, using
  /// `DC` for lookups.
  void VisitComment(const clang::RawComment* Comment,
                    const clang::DeclContext* DC,
                    const GraphObserver::NodeId& DocumentedNode);

  /// \brief Emit data for attributes attached to `Decl`, whose `NodeId`
  /// is `TargetNode`.
  void VisitAttributes(const clang::Decl* Decl,
                       const GraphObserver::NodeId& TargetNode);

  /// \brief Attempts to find the ID of the first parent of `Decl` for
  /// attaching a `childof` relationship.
  absl::optional<GraphObserver::NodeId> GetDeclChildOf(const clang::Decl* D);

  /// \brief Attempts to add some representation of `ND` to `Ostream`.
  /// \return true on success; false on failure.
  bool AddNameToStream(llvm::raw_string_ostream& Ostream,
                       const clang::NamedDecl* ND);

  /// \brief Assign `ND` (whose node ID is `TargetNode`) a USR if USRs are
  /// enabled.
  ///
  /// USRs are added only for NamedDecls that:
  ///   * are not under implicit template instantiations
  ///   * are not in a DeclContext inside a function body
  ///   * can actually be assigned USRs from Clang
  ///
  /// Similar to the way we deal with JVM names, the corpus, path,
  /// and root fields of a usr vname are cleared. Clients are permitted
  /// to write their own USR tickets. The USR value itself is encoded
  /// in capital hex (to match Clang's own internal USR stringification,
  /// modulo the configurable size of the SHA1 prefix).
  void AssignUSR(const GraphObserver::NodeId& TargetNode,
                 const clang::NamedDecl* ND);

  /// Assigns a USR to an alias.
  void AssignUSR(const GraphObserver::NameId& TargetName,
                 const GraphObserver::NodeId& AliasedType,
                 const clang::NamedDecl* ND) {
    AssignUSR(Observer.nodeIdForTypeAliasNode(TargetName, AliasedType), ND);
  }

  GraphObserver::NodeId ApplyBuiltinTypeConstructor(
      const char* BuiltinName, const GraphObserver::NodeId& Param);

  /// \brief Ascribes a type to `AscribeTo`.
  /// \param Type The `TypeLoc` referring to the type
  /// \param DType A possibly deduced type (or simply Type->getType()).
  /// \param AscribeTo The node to which the type should be ascribed.
  ///
  /// `auto` does not update TypeSourceInfo records after deduction, so
  /// a deduced `auto` in the source text will appear to be undeduced.
  /// In this case, it's useful to query the object being ascribed for its
  /// unlocated QualType, as this does get updated.
  void AscribeSpelledType(const clang::TypeLoc& Type,
                          const clang::QualType& TrueType,
                          const GraphObserver::NodeId& AscribeTo);

  /// \brief Returns the parent of the given node, along with the index
  /// at which the node appears underneath each parent.
  ///
  /// Note that this will lazily compute the parents of all nodes
  /// and store them for later retrieval. Thus, the first call is O(n)
  /// in the number of AST nodes.
  ///
  /// Caveats and FIXMEs:
  /// Calculating the parent map over all AST nodes will need to load the
  /// full AST. This can be undesirable in the case where the full AST is
  /// expensive to create (for example, when using precompiled header
  /// preambles). Thus, there are good opportunities for optimization here.
  /// One idea is to walk the given node downwards, looking for references
  /// to declaration contexts - once a declaration context is found, compute
  /// the parent map for the declaration context; if that can satisfy the
  /// request, loading the whole AST can be avoided. Note that this is made
  /// more complex by statements in templates having multiple parents - those
  /// problems can be solved by building closure over the templated parts of
  /// the AST, which also avoids touching large parts of the AST.
  /// Additionally, we will want to add an interface to already give a hint
  /// where to search for the parents, for example when looking at a statement
  /// inside a certain function.
  ///
  /// 'NodeT' can be one of Decl, Stmt, Type, TypeLoc,
  /// NestedNameSpecifier or NestedNameSpecifierLoc.
  template <typename NodeT>
  const IndexedParent* getIndexedParent(const NodeT& Node) {
    return getIndexedParent(clang::ast_type_traits::DynTypedNode::create(Node));
  }

  /// \return true if `Decl` and all of the nodes underneath it are prunable.
  ///
  /// A subtree is prunable if it's "the same" in all possible indexer runs.
  /// This excludes, for example, certain template instantiations.
  bool declDominatesPrunableSubtree(const clang::Decl* Decl);

  const IndexedParent* getIndexedParent(
      const clang::ast_type_traits::DynTypedNode& Node);

  /// Initializes AllParents, if necessary, and then returns a pointer to it.
  const IndexedParentMap* getAllParents();

  /// A map from memoizable DynTypedNodes to their parent nodes
  /// and their child indices with respect to those parents.
  /// Filled on the first call to `getIndexedParents`.
  std::unique_ptr<IndexedParentMap> AllParents;

  /// Records information about the template `Template` wrapping the node
  /// `BodyId`, including the edge linking the template and its body. Returns
  /// the `NodeId` for the dominating template.
  template <typename TemplateDeclish>
  GraphObserver::NodeId RecordTemplate(
      const TemplateDeclish* Decl, const GraphObserver::NodeId& BodyDeclNode);

  /// Records information about the generic class by wrapping the node
  /// `BodyId`. Returns the `NodeId` for the dominating generic type.
  GraphObserver::NodeId RecordGenericClass(
      const clang::ObjCInterfaceDecl* IDecl,
      const clang::ObjCTypeParamList* TPL, const GraphObserver::NodeId& BodyId);

  /// \brief Returns a vector of NodeId for each template argument.
  absl::optional<std::vector<GraphObserver::NodeId>> BuildTemplateArgumentList(
      llvm::ArrayRef<clang::TemplateArgument> Args);
  absl::optional<std::vector<GraphObserver::NodeId>> BuildTemplateArgumentList(
      llvm::ArrayRef<clang::TemplateArgumentLoc> Args);

  /// Dumps information about `TypeContext` to standard error when looking for
  /// an entry at (`Depth`, `Index`).
  void DumpTypeContext(unsigned Depth, unsigned Index);

  /// \brief Attempts to add a childof edge between DeclNode and its parent.
  /// \param Decl The (outer, in the case of a template) decl.
  /// \param DeclNode The (outer) `NodeId` for `Decl`.
  void AddChildOfEdgeToDeclContext(const clang::Decl* Decl,
                                   const GraphObserver::NodeId& DeclNode);

  /// Points at the inner node of the DeclContext, if it's a template.
  /// Otherwise points at the DeclContext as a Decl.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForDeclContext(
      const clang::DeclContext* DC);

  /// Points at the tapp node for a DeclContext, if it's an implicit template
  /// instantiation. Otherwise behaves as `BuildNodeIdForDeclContext`.
  absl::optional<GraphObserver::NodeId> BuildNodeIdForRefToDeclContext(
      const clang::DeclContext* DC);

  /// Avoid regenerating type node IDs and keep track of where we are while
  /// generating node IDs for recursive types. The key is opaque and
  /// makes sense only within the implementation of this class.
  TypeMap<NodeSet> TypeNodes;

  /// \brief Visit an Expr that refers to some NamedDecl.
  ///
  /// DeclRefExpr and ObjCIvarRefExpr are similar entities and can be processed
  /// in the same way but do not have a useful common ancestry.
  ///
  /// \param IsInit set to true if this is an initializing reference.
  bool VisitDeclRefOrIvarRefExpr(const clang::Expr* Expr,
                                 const clang::NamedDecl* const FoundDecl,
                                 clang::SourceLocation SL,
                                 bool IsImplicit = false, bool IsInit = false);

  /// \brief Connect a NodeId to the super and implemented protocols for a
  /// ObjCInterfaceDecl.
  ///
  /// Helper method used to connect an interface to the super and protocols it
  /// implements. It is also used to connect the interface implementation to
  /// these nodes as well. In that case, the interface implementation node is
  /// passed in as the first argument and the interface decl is passed in as the
  /// second.
  ///
  /// \param BodyDeclNode The node to connect to the super and protocols for
  /// the interface. This may be a interface decl node or an interface impl
  /// node.
  /// \param IFace The interface decl to use to look up the super and
  //// protocols.
  void ConnectToSuperClassAndProtocols(const GraphObserver::NodeId BodyDeclNode,
                                       const clang::ObjCInterfaceDecl* IFace);
  void ConnectToProtocols(const GraphObserver::NodeId BodyDeclNode,
                          clang::ObjCProtocolList::loc_iterator locStart,
                          clang::ObjCProtocolList::loc_iterator locEnd,
                          clang::ObjCProtocolList::iterator itStart,
                          clang::ObjCProtocolList::iterator itEnd);

  /// \brief Connect a parameter to a function decl.
  ///
  /// For FunctionDecl and ObjCMethodDecl, this connects the parameters to the
  /// function/method decl.
  /// \param Decl This should be a FunctionDecl or ObjCMethodDecl.
  void ConnectParam(const clang::Decl* Decl,
                    const GraphObserver::NodeId& FuncNode,
                    bool IsFunctionDefinition, const unsigned int ParamNumber,
                    const clang::ParmVarDecl* Param, bool DeclIsImplicit);

  /// \brief Draw the completes edge from a Decl to each of its redecls.
  void RecordCompletesForRedecls(const clang::Decl* Decl,
                                 const clang::SourceRange& NameRange,
                                 const GraphObserver::NodeId& DeclNode);

  /// \brief Draw an extends/category edge from the category to the class the
  /// category is extending.
  ///
  /// For example, @interface A (Cat) ... We draw an extends edge from the
  /// ObjCCategoryDecl for Cat to the ObjCInterfaceDecl for A.
  ///
  /// \param DeclNode The node for the category (impl or decl).
  /// \param IFace The class interface for class we are adding a category to.
  void ConnectCategoryToBaseClass(const GraphObserver::NodeId& DeclNode,
                                  const clang::ObjCInterfaceDecl* IFace);

  void LogErrorWithASTDump(const std::string& msg,
                           const clang::Decl* Decl) const;
  void LogErrorWithASTDump(const std::string& msg,
                           const clang::Expr* Expr) const;

  /// \brief This is used to handle the visitation of a clang::TypedefDecl
  /// or a clang::TypeAliasDecl.
  bool VisitCTypedef(const clang::TypedefNameDecl* Decl);

  /// \brief Find the implementation for `MD`. If `MD` is a definition, `MD` is
  /// returned. Otherwise, the method tries to find the implementation by
  /// looking through the interface and its implementation. If a method
  /// implementation is found, it is returned otherwise `MD` is returned.
  const clang::ObjCMethodDecl* FindMethodDefn(
      const clang::ObjCMethodDecl* MD, const clang::ObjCInterfaceDecl* I);

  void VisitObjCInterfaceDeclComment(
      const clang::ObjCInterfaceDecl* Decl, const clang::RawComment* Comment,
      const clang::DeclContext* DCxt,
      absl::optional<GraphObserver::NodeId> DCID);

  void VisitRecordDeclComment(const clang::RecordDecl* Decl,
                              const clang::RawComment* Comment,
                              const clang::DeclContext* DCxt,
                              absl::optional<GraphObserver::NodeId> DCID);

  /// \brief Maps known Decls to their NodeIds.
  llvm::DenseMap<const clang::Decl*, GraphObserver::NodeId> DeclToNodeId;

  /// \brief Used for calculating semantic hashes.
  SemanticHash Hash{
      [this](const clang::Decl* Decl) {
        return BuildNameIdForDecl(Decl).ToString();
      },
      // These enums are intentionally compatible.
      static_cast<SemanticHash::OnUnimplemented>(IgnoreUnimplemented)};

  /// \brief Enabled library-specific callbacks.
  const LibrarySupports& Supports;

  /// \brief The `Sema` instance to use.
  clang::Sema& Sema;

  /// \brief The cache to use to generate signatures.
  MarkedSourceCache MarkedSources;

  /// \return true if we should stop indexing.
  std::function<bool()> ShouldStopIndexing;

  /// \brief The active indexing job.
  std::unique_ptr<IndexJob> Job;

  /// \brief The controlling worklist.
  std::unique_ptr<IndexerWorklist> Worklist;

  /// \brief Comments we've already visited.
  std::unordered_set<const clang::RawComment*> VisitedComments;

  /// \brief The number of (raw) bytes to use to represent a USR. If 0,
  /// no USRs will be recorded.
  int UsrByteSize = 0;
};

/// \brief An `ASTConsumer` that passes events to a `GraphObserver`.
class IndexerASTConsumer : public clang::SemaConsumer {
 public:
  explicit IndexerASTConsumer(
      GraphObserver* GO, BehaviorOnUnimplemented B, BehaviorOnTemplates T,
      Verbosity V, BehaviorOnFwdDeclComments ObjC,
      BehaviorOnFwdDeclComments Cpp, const LibrarySupports& S,
      std::function<bool()> ShouldStopIndexing,
      std::function<std::unique_ptr<IndexerWorklist>(IndexerASTVisitor*)>
          CreateWorklist,
      int UsrByteSize)
      : Observer(GO),
        IgnoreUnimplemented(B),
        TemplateMode(T),
        Verbosity(V),
        ObjCFwdDocs(ObjC),
        CppFwdDocs(Cpp),
        Supports(S),
        ShouldStopIndexing(std::move(ShouldStopIndexing)),
        CreateWorklist(std::move(CreateWorklist)),
        UsrByteSize(UsrByteSize) {}

  void HandleTranslationUnit(clang::ASTContext& Context) override {
    CHECK(Sema != nullptr);
    IndexerASTVisitor Visitor(Context, IgnoreUnimplemented, TemplateMode,
                              Verbosity, ObjCFwdDocs, CppFwdDocs, Supports,
                              *Sema, ShouldStopIndexing, Observer, UsrByteSize);
    {
      ProfileBlock block(Observer->getProfilingCallback(), "traverse_tu");
      Visitor.Work(Context.getTranslationUnitDecl(), CreateWorklist(&Visitor));
    }
  }

  void InitializeSema(clang::Sema& S) override { Sema = &S; }

  void ForgetSema() override { Sema = nullptr; }

 private:
  GraphObserver* const Observer;
  /// Whether we should stop on missing cases or continue on.
  BehaviorOnUnimplemented IgnoreUnimplemented;
  /// Whether we should visit template instantiations.
  BehaviorOnTemplates TemplateMode;
  /// Whether we should emit all data.
  enum Verbosity Verbosity;
  /// Should we emit documentation for forward class decls in ObjC?
  BehaviorOnFwdDeclComments ObjCFwdDocs;
  /// Should we emit documentation for forward decls in C++?
  BehaviorOnFwdDeclComments CppFwdDocs;
  /// Which library supports are enabled.
  const LibrarySupports& Supports;
  /// The active Sema instance.
  clang::Sema* Sema;
  /// \return true if we should stop indexing.
  std::function<bool()> ShouldStopIndexing;
  /// \return a new worklist for the given visitor.
  std::function<std::unique_ptr<IndexerWorklist>(IndexerASTVisitor*)>
      CreateWorklist;
  /// \brief The number of (raw) bytes to use to represent a USR. If 0,
  /// no USRs will be recorded.
  int UsrByteSize = 0;
};

}  // namespace kythe

#endif  // KYTHE_CXX_INDEXER_CXX_INDEXER_AST_HOOKS_H_
