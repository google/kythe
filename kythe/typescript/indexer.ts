/*
 * Copyright 2017 The Kythe Authors. All rights reserved.
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

import 'source-map-support/register';
import * as fs from 'fs';
import * as path from 'path';
import * as ts from 'typescript';

import * as utf8 from './utf8';

/** VName is the type of Kythe node identities. */
export interface VName {
  signature: string;
  corpus: string;
  root: string;
  path: string;
  language: string;
}

/**
 * toArray converts an Iterator to an array of its values.
 * It's necessary when running in ES5 environments where for-of loops
 * don't iterate through Iterators.
 */
function toArray<T>(it: Iterator<T>): T[] {
  const array: T[] = [];
  for (let next = it.next(); !next.done; next = it.next()) {
    array.push(next.value);
  }
  return array;
}

/**
 * stripExtension strips the .d.ts or .ts extension from a path.
 * It's used to map a file path to the module name.
 */
function stripExtension(path: string): string {
  return path.replace(/\.(d\.)?ts$/, '');
}

/**
 * TSNamespace represents the two namespaces of TypeScript: types and values.
 * A given symbol may be a type, it may be a value, and the two may even
 * be unrelated.
 *
 * See the table at
 *   https://www.typescriptlang.org/docs/handbook/declaration-merging.html
 *
 * TODO: there are actually three namespaces; the third is (confusingly)
 * itself called namespaces.  Implement those in this enum and other places.
 */
enum TSNamespace {
  TYPE,
  VALUE,
}

/** Visitor manages the indexing process for a single TypeScript SourceFile. */
class Vistor {
  /** kFile is the VName for the 'file' node representing the source file. */
  kFile: VName;

  /**
   * symbolNames maps ts.Symbols to their assigned VNames.
   * The value is a tuple of the separate TypeScript namespaces, and entries
   * in it correspond to TSNamespace values.  See the documentation of
   * TSNamespace.
   */
  symbolNames = new Map<ts.Symbol, [VName | null, VName|null]>();

  /**
   * anonId increments for each anonymous block, to give them unique
   * signatures.
   */
  anonId = 0;

  /**
   * anonNames maps nodes to the anonymous names assigned to them.
   */
  anonNames = new Map<ts.Node, string>();

  typeChecker: ts.TypeChecker;

  /** A shorter name for the rootDir in the CompilerOptions. */
  sourceRoot: string;

  /**
   * rootDirs is the list of rootDirs in the compiler options, sorted
   * longest first.  See this.moduleName().
   */
  rootDirs: string[];

  constructor(
      /**
       * The VName for the CompilationUnit, containing compilation-wide info.
       */
      private compilationUnit: VName,
      /**
       * A map of path to path-specific VName.
       */
      private pathVNames: Map<string, VName>, program: ts.Program,
      private file: ts.SourceFile,
      private getOffsetTable: (path: string) => utf8.OffsetTable) {
    this.typeChecker = program.getTypeChecker();

    this.sourceRoot = program.getCompilerOptions().rootDir || process.cwd();
    let rootDirs = program.getCompilerOptions().rootDirs || [this.sourceRoot];
    rootDirs = rootDirs.map(d => d + '/');
    rootDirs.sort((a, b) => b.length - a.length);
    this.rootDirs = rootDirs;

    this.kFile = this.newVName(
        /* empty signature */ '',
        path.relative(this.sourceRoot, file.fileName));
    this.kFile.language = '';
  }

  /**
   * emit emits a Kythe entry, structured as a JSON object.  Defaults to
   * emitting to stdout but users may replace it.
   */
  emit = (obj: {}) => {
    console.log(JSON.stringify(obj));
  };

  todo(node: ts.Node, message: string) {
    const sourceFile = node.getSourceFile();
    const file = path.relative(this.sourceRoot, sourceFile.fileName);
    const {line, character} =
        ts.getLineAndCharacterOfPosition(sourceFile, node.getStart());
    console.warn(`TODO: ${file}:${line}:${character}: ${message}`);
  }

  /**
   * moduleName returns the ES6 module name of a path to a source file.
   * E.g. foo/bar.ts and foo/bar.d.ts both have the same module name,
   * 'foo/bar', and rootDirs (like bazel-bin/) are eliminated.
   * See README.md for a discussion of this.
   */
  moduleName(sourcePath: string): string {
    // Compute sourcePath as relative to one of the rootDirs.
    // This canonicalizes e.g. bazel-bin/foo to just foo.
    // Note that this.rootDirs is sorted longest first, so we'll use the
    // longest match.
    for (const rootDir of this.rootDirs) {
      if (sourcePath.startsWith(rootDir)) {
        sourcePath = path.relative(rootDir, sourcePath);
        break;
      }
    }
    return stripExtension(sourcePath);
  }

  /**
   * newVName returns a new VName with a given signature and path.
   */
  newVName(signature: string, path: string): VName {
    const vname = this.pathVNames.get(path);
    return {
      signature,
      corpus: vname && vname.corpus ? vname.corpus :
                                      this.compilationUnit.corpus,
      root: vname && vname.root ? vname.root : this.compilationUnit.root,
      path: vname && vname.path ? vname.path : path,
      language: 'typescript',
    };
  }

  /** newAnchor emits a new anchor entry that covers a TypeScript node. */
  newAnchor(node: ts.Node, start = node.getStart(), end = node.end): VName {
    const name = this.newVName(
        `@${start}:${end}`,
        // An anchor refers to specific text, so its path is the file path,
        // not the module name.
        path.relative(this.sourceRoot, node.getSourceFile().fileName));
    this.emitNode(name, 'anchor');
    const offsetTable = this.getOffsetTable(node.getSourceFile().fileName);
    this.emitFact(name, 'loc/start', offsetTable.lookup(start).toString());
    this.emitFact(name, 'loc/end', offsetTable.lookup(end).toString());
    this.emitEdge(name, 'childof', this.kFile);
    return name;
  }

  /** emitNode emits a new node entry, declaring the kind of a VName. */
  emitNode(source: VName, kind: string) {
    this.emitFact(source, 'node/kind', kind);
  }

  /** emitFact emits a new fact entry, tying an attribute to a VName. */
  emitFact(source: VName, name: string, value: string) {
    this.emit({
      source,
      fact_name: '/kythe/' + name,
      fact_value: Buffer.from(value).toString('base64'),
    });
  }

  /** emitEdge emits a new edge entry, relating two VNames. */
  emitEdge(source: VName, name: string, target: VName) {
    this.emit({
      source,
      edge_kind: '/kythe/edge/' + name,
      target,
      fact_name: '/',
    });
  }

  /**
   * anonName assigns a freshly generated name to a Node.
   * It's used to give stable names to e.g. anonymous objects.
   */
  anonName(node: ts.Node): string {
    let name = this.anonNames.get(node);
    if (!name) {
      name = `anon${this.anonId++}`;
      this.anonNames.set(node, name);
    }
    return name;
  }

  /**
   * scopedSignature computes a scoped name for a ts.Node.
   * E.g. if you have a function `foo` containing a block containing a variable
   * `bar`, it might return a VName like
   *   signature: "foo.block0.bar""
   *   path: <appropriate path to module>
   */
  scopedSignature(startNode: ts.Node): VName {
    let moduleName: string|undefined;
    const parts: string[] = [];

    // Traverse the containing blocks upward, gathering names from nodes that
    // introduce scopes.
    for (let node: ts.Node|undefined = startNode; node != null;
         node = node.parent) {
      switch (node.kind) {
        case ts.SyntaxKind.ExportAssignment:
          const exportDecl = node as ts.ExportAssignment;
          if (!exportDecl.isExportEquals) {
            // It's an "export default" statement.
            // This is semantically equivalent to exporting a variable
            // named 'default'.
            parts.push('default');
          } else {
            this.todo(node, 'handle ExportAssignment with =');
          }
          break;
        case ts.SyntaxKind.ArrowFunction:
          // Arrow functions are anonymous, so generate a unique id.
          parts.push(`arrow${this.anonId++}`);
          break;
        case ts.SyntaxKind.Block:
          if (node.parent &&
              (node.parent.kind === ts.SyntaxKind.FunctionDeclaration ||
               node.parent.kind === ts.SyntaxKind.MethodDeclaration)) {
            // A block that's an immediate child of a function is the
            // function's body, so it doesn't need a separate name.
            continue;
          }
          parts.push(`block${this.anonId++}`);
          break;
        case ts.SyntaxKind.BindingElement:
        case ts.SyntaxKind.ClassDeclaration:
        case ts.SyntaxKind.ClassExpression:
        case ts.SyntaxKind.EnumDeclaration:
        case ts.SyntaxKind.EnumMember:
        case ts.SyntaxKind.FunctionDeclaration:
        case ts.SyntaxKind.InterfaceDeclaration:
        case ts.SyntaxKind.ImportSpecifier:
        case ts.SyntaxKind.ExportSpecifier:
        case ts.SyntaxKind.MethodDeclaration:
        case ts.SyntaxKind.MethodSignature:
        case ts.SyntaxKind.NamespaceImport:
        case ts.SyntaxKind.ObjectLiteralExpression:
        case ts.SyntaxKind.Parameter:
        case ts.SyntaxKind.PropertyAssignment:
        case ts.SyntaxKind.PropertyDeclaration:
        case ts.SyntaxKind.PropertySignature:
        case ts.SyntaxKind.TypeAliasDeclaration:
        case ts.SyntaxKind.TypeParameter:
        case ts.SyntaxKind.VariableDeclaration:
          const decl = node as ts.NamedDeclaration;
          if (decl.name && decl.name.kind === ts.SyntaxKind.Identifier) {
            parts.push(decl.name.text);
          } else {
            // TODO: handle other declarations, e.g. binding patterns.
            parts.push(this.anonName(node));
          }
          break;
        case ts.SyntaxKind.Constructor:
          parts.push('constructor');
          break;
        case ts.SyntaxKind.ModuleDeclaration:
          const modDecl = node as ts.ModuleDeclaration;
          if (modDecl.name.kind === ts.SyntaxKind.StringLiteral) {
            // Syntax like:
            //   declare module 'foo/bar' {}
            // This is the syntax for defining symbols in another, named
            // module.
            moduleName = (modDecl.name as ts.StringLiteral).text;
          } else if (modDecl.name.kind === ts.SyntaxKind.Identifier) {
            // Syntax like:
            //   declare module foo {}
            // without quotes is just an obsolete way of saying 'namespace'.
            parts.push((modDecl.name as ts.Identifier).text);
          }
          break;
        case ts.SyntaxKind.SourceFile:
          // moduleName can already be set if the target was contained within
          // a "declare module 'foo/bar'" block (see the handling of
          // ModuleDeclaration).  Otherwise, the module name is derived from the
          // name of the current file.
          if (!moduleName) {
            moduleName = this.moduleName((node as ts.SourceFile).fileName);
          }
          break;
        default:
          // Most nodes are children of other nodes that do not introduce a
          // new namespace, e.g. "return x;", so ignore all other parents
          // by default.
          // TODO: namespace {}, etc.

          // If the node is actually some subtype that has a 'name' attribute
          // it's likely this function should have handled it.  Dynamically
          // probe for this case and warn if we missed one.
          if ('name' in (node as any)) {
            this.todo(
                node,
                `scopedSignature: ${ts.SyntaxKind[node.kind]} ` +
                    `has unused 'name' property`);
          }
      }
    }

    // The names were gathered from bottom to top, so reverse before joining.
    const sig = parts.reverse().join('.');
    return this.newVName(sig, moduleName!);
  }

  /**
   * getSymbolAtLocation is the same as this.typeChecker.getSymbolAtLocation,
   * except that it has a return type that properly captures that
   * getSymbolAtLocation can return undefined.  (The TypeScript API itself is
   * not yet null-safe, so it hasn't been annotated with full types.)
   */
  getSymbolAtLocation(node: ts.Node): ts.Symbol|undefined {
    return this.typeChecker.getSymbolAtLocation(node);
  }

  /** getSymbolName computes the VName (and signature) of a ts.Symbol. */
  getSymbolName(sym: ts.Symbol, ns: TSNamespace): VName {
    let vnames = this.symbolNames.get(sym);
    if (vnames && vnames[ns]) return vnames[ns]!;

    if (!sym.declarations || sym.declarations.length < 1) {
      throw new Error('TODO: symbol has no declarations?');
    }
    // TODO: think about symbols with multiple declarations.

    const decl = sym.declarations[0];
    const vname = this.scopedSignature(decl);
    // The signature of a value is undecorated;
    // the signature of a type has the #type suffix.
    if (ns === TSNamespace.TYPE) {
      vname.signature += '#type';
    }

    // Save it in the appropriate slot in the symbolNames table.
    if (!vnames) vnames = [null, null];
    vnames[ns] = vname;
    this.symbolNames.set(sym, vnames);

    return vname;
  }

  visitTypeParameters(params: ReadonlyArray<ts.TypeParameterDeclaration>) {
    for (const param of params) {
      const sym = this.getSymbolAtLocation(param.name);
      if (!sym) {
        this.todo(param, `type param ${param.getText()} has no symbol`);
        return;
      }
      const kType = this.getSymbolName(sym, TSNamespace.TYPE);
      this.emitNode(kType, 'absvar');
      this.emitEdge(this.newAnchor(param.name), 'defines/binding', kType);
    }
  }

  /**
   * visitHeritage visits the X found in an 'extends X' or 'implements X'.
   *
   * These are subtle in an interesting way.  When you have
   *   interface X extends Y {}
   * that is referring to the *type* Y (because interfaces are types, not
   * values).  But it's also legal to write
   *   class X extends (class Z { ... }) {}
   * where the thing in the extends clause is itself an expression, and the
   * existing logic for visiting a class expression already handles modelling
   * the class as both a type and a value.
   *
   * The full set of possible combinations is:
   * - class extends => value
   * - interface extends => type
   * - class implements => type
   * - interface implements => illegal
   */
  visitHeritage(heritageClauses: ReadonlyArray<ts.HeritageClause>) {
    for (const heritage of heritageClauses) {
      if (heritage.token === ts.SyntaxKind.ExtendsKeyword && heritage.parent &&
          heritage.parent.kind !== ts.SyntaxKind.InterfaceDeclaration) {
        this.visit(heritage);
      } else {
        this.visitType(heritage);
      }
    }
  }

  visitInterfaceDeclaration(decl: ts.InterfaceDeclaration) {
    const sym = this.getSymbolAtLocation(decl.name);
    if (!sym) {
      this.todo(decl.name, `interface ${decl.name.getText()} has no symbol`);
      return;
    }
    const kType = this.getSymbolName(sym, TSNamespace.TYPE);
    this.emitNode(kType, 'interface');
    this.emitEdge(this.newAnchor(decl.name), 'defines/binding', kType);

    if (decl.typeParameters) this.visitTypeParameters(decl.typeParameters);
    if (decl.heritageClauses) this.visitHeritage(decl.heritageClauses);
    for (const member of decl.members) {
      this.visit(member);
    }
  }

  visitTypeAliasDeclaration(decl: ts.TypeAliasDeclaration) {
    const sym = this.getSymbolAtLocation(decl.name);
    if (!sym) {
      this.todo(decl.name, `type alias ${decl.name.getText()} has no symbol`);
      return;
    }
    const kType = this.getSymbolName(sym, TSNamespace.TYPE);
    this.emitNode(kType, 'talias');
    this.emitEdge(this.newAnchor(decl.name), 'defines/binding', kType);

    if (decl.typeParameters) this.visitTypeParameters(decl.typeParameters);
    this.visitType(decl.type);
    // TODO: in principle right here we emit an "aliases" edge.
    // However, it's complicated by the fact that some types don't have
    // specific names to alias, e.g.
    //   type foo = number|string;
    // Just punt for now.
  }

  /**
   * visitType is the main dispatch for visiting type nodes.
   * It's separate from visit() because bare ts.Identifiers within a normal
   * expression are values (handled by visit) but bare ts.Identifiers within
   * a type are types (handled here).
   */
  visitType(node: ts.Node): void {
    switch (node.kind) {
      case ts.SyntaxKind.Identifier:
        const sym = this.getSymbolAtLocation(node);
        if (!sym) {
          this.todo(node, `type ${node.getText()} has no symbol`);
          return;
        }
        const name = this.getSymbolName(sym, TSNamespace.TYPE);
        this.emitEdge(this.newAnchor(node), 'ref', name);
        return;
      default:
        // Default recursion, but using visitType(), not visit().
        return ts.forEachChild(node, n => this.visitType(n));
    }
  }

  /**
   * getPathFromModule gets the "module path" from the module import
   * symbol referencing a module system path to reference to a module.
   *
   * E.g. from
   *   import ... from './foo';
   * getPathFromModule(the './foo' node) might return a string like
   * 'path/to/project/foo'.  See this.moduleName().
   */
  getModulePathFromModuleReference(sym: ts.Symbol): string {
    const name = sym.name;
    // TODO: this is hacky; it may be the case we need to use the TypeScript
    // module resolver to get the real path (?).  But it appears the symbol
    // name is the quoted(!) path to the module.
    if (!(name.startsWith('"') && name.endsWith('"'))) {
      throw new Error(`TODO: handle module symbol ${name}`);
    }
    const sourcePath = name.substr(1, name.length - 2);
    return this.moduleName(sourcePath);
  }

  /**
   * visitImportSpecifier handles a single entry in an import statement, e.g.
   * "bar" in code like
   *   import {foo, bar} from 'baz';
   * See visitImportDeclaration for the code handling the entire statement.
   *
   * @return The VName for the import.
   */
  visitImport(name: ts.Node): VName {
    // An import both aliases a symbol from another module
    // (call it the "remote" symbol) and it defines a local symbol.
    //
    // Those two symbols often have the same name, with statements like:
    //   import {foo} from 'bar';
    // But they can be different, in e.g.
    //   import {foo as baz} from 'bar';
    // Which maps the remote symbol named 'foo' to a local named 'baz'.
    //
    // In all cases TypeScript maintains two different ts.Symbol objects,
    // one for the local and one for the remote.  In principle for the
    // simple import statement
    //   import {foo} from 'bar';
    // the "foo" should be both:
    // - a ref/imports to the remote symbol
    // - a defines/binding for the local symbol
    //
    // But in Kythe the UI for stacking multiple references out from a single
    // anchor isn't great, so this code instead unifies all references
    // (including renaming imports) to a single VName.

    const localSym = this.getSymbolAtLocation(name);
    if (!localSym) {
      throw new Error(`TODO: local name ${name} has no symbol`);
    }

    const remoteSym = this.typeChecker.getAliasedSymbol(localSym);
    // This imported symbol can refer to a type, a value, or both.
    const kImportValue = remoteSym.flags & ts.SymbolFlags.Value ?
        this.getSymbolName(remoteSym, TSNamespace.VALUE) :
        null;
    const kImportType = remoteSym.flags & ts.SymbolFlags.Type ?
        this.getSymbolName(remoteSym, TSNamespace.TYPE) :
        null;
    // Mark the local symbol with the remote symbol's VName so that all
    // references resolve to the remote symbol.
    this.symbolNames.set(localSym, [kImportType, kImportValue]);

    // The name anchor must link somewhere.  In rare cases a symbol is both
    // a type and a value that resolve to two different locations; for now,
    // because we must choose one, just prefer linking to the value.
    // One of the value or type reference should be non-null.
    const kImport = (kImportValue || kImportType)!;
    this.emitEdge(this.newAnchor(name), 'ref/imports', kImport);
    return kImport;
  }

  /** visitImportDeclaration handles the various forms of "import ...". */
  visitImportDeclaration(decl: ts.ImportDeclaration) {
    // All varieties of import statements reference a module on the right,
    // so start by linking that.
    const moduleSym = this.getSymbolAtLocation(decl.moduleSpecifier);
    if (!moduleSym) {
      // This can occur when the module failed to resolve to anything.
      // See testdata/import_missing.ts for more on how that could happen.
      return;
    }
    const kModule = this.newVName(
        'module', this.getModulePathFromModuleReference(moduleSym));
    this.emitEdge(this.newAnchor(decl.moduleSpecifier), 'ref/imports', kModule);

    if (!decl.importClause) {
      // This is a side-effecting import that doesn't declare anything, e.g.:
      //   import 'foo';
      return;
    }
    const clause = decl.importClause;

    if (clause.name) {
      // This is a default import, e.g.:
      //   import foo from './bar';
      this.visitImport(clause.name);
      return;
    }

    if (!clause.namedBindings) {
      // TODO: I believe clause.name or clause.namedBindings are always present,
      // which means this check is not necessary, but the types don't show that.
      throw new Error(`import declaration ${decl.getText()} has no bindings`);
    }
    switch (clause.namedBindings.kind) {
      case ts.SyntaxKind.NamespaceImport:
        // This is a namespace import, e.g.:
        //   import * as foo from 'foo';
        const name = clause.namedBindings.name;
        const sym = this.getSymbolAtLocation(name);
        if (!sym) {
          this.todo(
              clause, `namespace import ${clause.getText()} has no symbol`);
          return;
        }
        const kModuleObject = this.getSymbolName(sym, TSNamespace.VALUE);
        this.emitNode(kModuleObject, 'variable');
        this.emitEdge(this.newAnchor(name), 'defines/binding', kModuleObject);
        break;
      case ts.SyntaxKind.NamedImports:
        // This is named imports, e.g.:
        //   import {bar, baz} from 'foo';
        const imports = clause.namedBindings.elements;
        for (const imp of imports) {
          const kImport = this.visitImport(imp.name);
          if (imp.propertyName) {
            this.emitEdge(
                this.newAnchor(imp.propertyName), 'ref/imports', kImport);
          }
        }
        break;
    }
  }

  /**
   * When a file imports another file, with syntax like
   *   import * as x from 'some/path';
   * we wants 'some/path' to refer to a VName that just means "the entire
   * file".  It doesn't refer to any text in particular, so we just mark
   * the first letter in the file as the anchor for this.
   */
  emitModuleAnchor(sf: ts.SourceFile) {
    const kMod = this.newVName('module', this.moduleName(this.file.fileName));
    this.emitFact(kMod, 'node/kind', 'record');
    this.emitEdge(this.kFile, 'childof', kMod);

    // Emit the anchor, bound to the beginning of the file.
    const anchor = this.newAnchor(this.file, 0, 1);
    this.emitEdge(anchor, 'defines/binding', kMod);
  }

  /**
   * Handles code like:
   *   export default ...;
   *   export = ...;
   */
  visitExportAssignment(assign: ts.ExportAssignment) {
    if (assign.isExportEquals) {
      this.todo(assign, `handle export = statement`);
    } else {
      // export default <expr>;
      // is the same as exporting the expression under the symbol named
      // "default".  But we don't have a nice name to link the symbol to!
      // So instead we link the keyword "default" itself to the VName.
      // The TypeScript AST does not expose the location of the 'default'
      // keyword so we just find it in the source text to link it.
      const ofs = assign.getText().indexOf('default');
      if (ofs < 0) throw new Error(`'export default' without 'default'?`);
      const start = assign.getStart() + ofs;
      const anchor = this.newAnchor(assign, start, start + 'default'.length);
      this.emitEdge(anchor, 'defines/binding', this.scopedSignature(assign));
    }
  }

  /**
   * Handles code that explicitly exports a symbol, like:
   *   export {Foo} from './bar';
   *
   * Note that export can also be a modifier on a declaration, like:
   *   export class Foo {}
   * and that case is handled as part of the ordinary declaration handling.
   */
  visitExportDeclaration(decl: ts.ExportDeclaration) {
    if (decl.exportClause) {
      for (const exp of decl.exportClause.elements) {
        const localSym = this.getSymbolAtLocation(exp.name);
        if (!localSym) {
          console.error(`TODO: export ${name} has no symbol`);
          continue;
        }
        // TODO: import a type, not just a value.
        const remoteSym = this.typeChecker.getAliasedSymbol(localSym);
        const kExport = this.getSymbolName(remoteSym, TSNamespace.VALUE);
        this.emitEdge(this.newAnchor(exp.name), 'ref', kExport);
        if (exp.propertyName) {
          // Aliased export; propertyName is the 'as <...>' bit.
          this.emitEdge(this.newAnchor(exp.propertyName), 'ref', kExport);
        }
      }
    }
    if (decl.moduleSpecifier) {
      this.todo(
          decl.moduleSpecifier, `handle module specifier in ${decl.getText()}`);
    }
  }

  visitVariableStatement(stmt: ts.VariableStatement) {
    // A VariableStatement contains potentially multiple variable declarations,
    // as in:
    //   var x = 3, y = 4;
    // In the (common) case where there's a single variable declared, we look
    // for documentation for that variable above the statement.
    let vname: VName|undefined;
    for (const decl of stmt.declarationList.declarations) {
      vname = this.visitVariableDeclaration(decl);
    }
    if (stmt.declarationList.declarations.length === 1 && vname !== undefined) {
      this.visitJSDoc(stmt, vname);
    }
  }

  /**
   * Note: visitVariableDeclaration is also used for class properties;
   * the decl parameter is the union of the attributes of the two types.
   * @return the generated VName for the declaration, if any.
   */
  visitVariableDeclaration(decl: {
    name: ts.BindingName|ts.PropertyName,
    type?: ts.TypeNode,
    initializer?: ts.Expression, kind: ts.SyntaxKind,
  }): VName|undefined {
    let vname: VName|undefined;
    switch (decl.name.kind) {
      case ts.SyntaxKind.Identifier:
        const sym = this.getSymbolAtLocation(decl.name);
        if (!sym) {
          this.todo(
              decl.name, `declaration ${decl.name.getText()} has no symbol`);
          return undefined;
        }
        vname = this.getSymbolName(sym, TSNamespace.VALUE);
        this.emitNode(vname, 'variable');

        this.emitEdge(this.newAnchor(decl.name), 'defines/binding', vname);
        break;
      case ts.SyntaxKind.ObjectBindingPattern:
      case ts.SyntaxKind.ArrayBindingPattern:
        for (const element of (decl.name as ts.BindingPattern).elements) {
          this.visit(element);
        }
        break;
      default:
        this.todo(
            decl.name,
            `handle variable declaration: ${ts.SyntaxKind[decl.name.kind]}`);
    }
    if (decl.type) this.visitType(decl.type);
    if (decl.initializer) this.visit(decl.initializer);
    if (vname && decl.kind === ts.SyntaxKind.PropertyDeclaration) {
      const declNode = decl as ts.PropertyDeclaration;
      if (ts.getCombinedModifierFlags(declNode) & ts.ModifierFlags.Static) {
        this.emitFact(vname, 'tag/static', '');
      }
    }
    return vname;
  }

  visitFunctionLikeDeclaration(decl: ts.FunctionLikeDeclaration) {
    this.visitDecorators(decl.decorators || []);
    let kFunc: VName|undefined = undefined;
    if (decl.name) {
      const sym = this.getSymbolAtLocation(decl.name);
      if (decl.name.kind === ts.SyntaxKind.ComputedPropertyName) {
        // TODO: it's not clear what to do with computed property named
        // functions.  They don't have a symbol.
        this.visit((decl.name as ts.ComputedPropertyName).expression);
      } else {
        if (!sym) {
          this.todo(
              decl.name,
              `function declaration ${decl.name.getText()} has no symbol`);
          return;
        }
        kFunc = this.getSymbolName(sym, TSNamespace.VALUE);
        this.emitNode(kFunc, 'function');

        this.emitEdge(this.newAnchor(decl.name), 'defines/binding', kFunc);

        this.visitJSDoc(decl, kFunc);
      }
    } else {
      // TODO: choose VName for anonymous functions.
      kFunc = this.newVName('TODO', 'TODOPath');
    }

    if (kFunc && decl.parent) {
      // Emit a "childof" edge on class/interface members.
      if (decl.parent.kind === ts.SyntaxKind.ClassDeclaration ||
          decl.parent.kind === ts.SyntaxKind.ClassExpression ||
          decl.parent.kind === ts.SyntaxKind.InterfaceDeclaration) {
        const parentName = (decl.parent as ts.ClassLikeDeclaration).name;
        if (parentName !== undefined) {
          const parentSym = this.getSymbolAtLocation(parentName);
          if (!parentSym) {
            this.todo(parentName, `parent ${parentName} has no symbol`);
            return;
          }
          const kParent = this.getSymbolName(parentSym, TSNamespace.TYPE);
          this.emitEdge(kFunc, 'childof', kParent);
        }
      }

      // TODO: emit an "overrides" edge if this method overrides.
      // It appears the TS API doesn't make finding that information easy,
      // perhaps because it mostly cares about whether types are structrually
      // compatible.  But I think you can start from the parent class/iface,
      // then from there follow the "implements"/"extends" chain to other
      // classes/ifaces, and then from there look for methods with matching
      // names.
    }

    for (const [index, param] of toArray(decl.parameters.entries())) {
      this.visitDecorators(param.decorators || []);
      const sym = this.getSymbolAtLocation(param.name);
      if (!sym) {
        this.todo(param.name, `param ${param.name.getText()} has no symbol`);
        continue;
      }
      const kParam = this.getSymbolName(sym, TSNamespace.VALUE);
      this.emitNode(kParam, 'variable');
      if (kFunc) this.emitEdge(kFunc, `param.${index}`, kParam);

      this.emitEdge(this.newAnchor(param.name), 'defines/binding', kParam);
      if (param.type) this.visitType(param.type);

      if (param.initializer) this.visit(param.initializer);
    }

    if (decl.type) {
      // "type" here is the return type of the function.
      this.visitType(decl.type);
    }

    if (decl.typeParameters) this.visitTypeParameters(decl.typeParameters);
    if (decl.body) {
      this.visit(decl.body);
    } else {
      if (kFunc) this.emitFact(kFunc, 'complete', 'incomplete');
    }
  }

  visitDecorators(decors: ReadonlyArray<ts.Decorator>) {
    for (const decor of decors) {
      this.visit(decor);
    }
  }

  visitClassDeclaration(decl: ts.ClassDeclaration) {
    this.visitDecorators(decl.decorators || []);
    if (decl.name) {
      const sym = this.getSymbolAtLocation(decl.name);
      if (!sym) {
        this.todo(decl.name, `class ${decl.name.getText()} has no symbol`);
        return;
      }
      // A 'class' declaration declares both a type (a 'record', representing
      // instances of the class) and a value (the constructor).
      const kClass = this.getSymbolName(sym, TSNamespace.TYPE);
      this.emitNode(kClass, 'record');
      const kClassCtor = this.getSymbolName(sym, TSNamespace.VALUE);
      this.emitNode(kClassCtor, 'function');
      // TODO: the specific constructor() should really be the thing tagged
      // with constructor, but we also need to handle classes that don't declare
      // a constructor.  Fix me later.
      this.emitFact(kClassCtor, 'subkind', 'constructor');

      const anchor = this.newAnchor(decl.name);
      this.emitEdge(anchor, 'defines/binding', kClass);
      this.emitEdge(anchor, 'defines/binding', kClassCtor);

      this.visitJSDoc(decl, kClass);
    }
    if (decl.typeParameters) this.visitTypeParameters(decl.typeParameters);
    if (decl.heritageClauses) this.visitHeritage(decl.heritageClauses);
    for (const member of decl.members) {
      this.visit(member);
    }
  }

  visitEnumDeclaration(decl: ts.EnumDeclaration) {
    const sym = this.getSymbolAtLocation(decl.name);
    if (!sym) return;
    const kType = this.getSymbolName(sym, TSNamespace.TYPE);
    this.emitNode(kType, 'record');
    const kValue = this.getSymbolName(sym, TSNamespace.VALUE);
    this.emitNode(kValue, 'constant');

    const anchor = this.newAnchor(decl.name);
    this.emitEdge(anchor, 'defines/binding', kType);
    this.emitEdge(anchor, 'defines/binding', kValue);
    for (const member of decl.members) {
      this.visit(member);
    }
  }

  visitEnumMember(decl: ts.EnumMember) {
    const sym = this.getSymbolAtLocation(decl.name);
    if (!sym) return;
    const kMember = this.getSymbolName(sym, TSNamespace.VALUE);
    this.emitNode(kMember, 'constant');
    this.emitEdge(this.newAnchor(decl.name), 'defines/binding', kMember);
  }

  /**
   * visitJSDoc attempts to attach a 'doc' node to a given target, by looking
   * for JSDoc comments.
   */
  visitJSDoc(node: ts.Node, target: VName) {
    const text = node.getFullText();
    const comments = ts.getLeadingCommentRanges(text, 0);
    if (!comments) return;

    let jsdoc: string|undefined;
    for (const commentRange of comments) {
      if (commentRange.kind !== ts.SyntaxKind.MultiLineCommentTrivia) continue;
      const comment =
          text.substring(commentRange.pos + 2, commentRange.end - 2);
      if (!comment.startsWith('*')) {
        // Not a JSDoc comment.
        continue;
      }
      // Strip the ' * ' bits that start lines within the comment.
      jsdoc = comment.replace(/^[ \t]*\* ?/mg, '');
      break;
    }
    if (jsdoc === undefined) return;

    // Strip leading and trailing whitespace.
    jsdoc = jsdoc.replace(/^\s+/, '').replace(/\s+$/, '');
    const doc = this.newVName(target.signature + '#doc', target.path);
    this.emitNode(doc, 'doc');
    this.emitEdge(doc, 'documents', target);
    this.emitFact(doc, 'text', jsdoc);
  }

  /** visit is the main dispatch for visiting AST nodes. */
  visit(node: ts.Node): void {
    switch (node.kind) {
      case ts.SyntaxKind.ImportDeclaration:
        return this.visitImportDeclaration(node as ts.ImportDeclaration);
      case ts.SyntaxKind.ExportAssignment:
        return this.visitExportAssignment(node as ts.ExportAssignment);
      case ts.SyntaxKind.ExportDeclaration:
        return this.visitExportDeclaration(node as ts.ExportDeclaration);
      case ts.SyntaxKind.VariableStatement:
        return this.visitVariableStatement(node as ts.VariableStatement);
      case ts.SyntaxKind.PropertyAssignment:  // property in object literal
      case ts.SyntaxKind.PropertyDeclaration:
      case ts.SyntaxKind.PropertySignature:
        const vname =
            this.visitVariableDeclaration(node as ts.PropertyDeclaration);
        if (vname) this.visitJSDoc(node, vname);
        return;
      case ts.SyntaxKind.ArrowFunction:
      case ts.SyntaxKind.Constructor:
      case ts.SyntaxKind.FunctionDeclaration:
      case ts.SyntaxKind.MethodDeclaration:
      case ts.SyntaxKind.MethodSignature:
        return this.visitFunctionLikeDeclaration(
            node as ts.FunctionLikeDeclaration);
      case ts.SyntaxKind.ClassDeclaration:
        return this.visitClassDeclaration(node as ts.ClassDeclaration);
      case ts.SyntaxKind.InterfaceDeclaration:
        return this.visitInterfaceDeclaration(node as ts.InterfaceDeclaration);
      case ts.SyntaxKind.TypeAliasDeclaration:
        return this.visitTypeAliasDeclaration(node as ts.TypeAliasDeclaration);
      case ts.SyntaxKind.EnumDeclaration:
        return this.visitEnumDeclaration(node as ts.EnumDeclaration);
      case ts.SyntaxKind.EnumMember:
        return this.visitEnumMember(node as ts.EnumMember);
      case ts.SyntaxKind.TypeReference:
        return this.visitType(node as ts.TypeNode);
      case ts.SyntaxKind.BindingElement:
        this.visitVariableDeclaration(node as ts.BindingElement);
        return;
      case ts.SyntaxKind.Identifier:
        // Assume that this identifer is occurring as part of an
        // expression; we handle identifiers that occur in other
        // circumstances (e.g. in a type) separately in visitType.
        const sym = this.getSymbolAtLocation(node);
        if (!sym) {
          // E.g. a field of an "any".
          return;
        }
        if (!sym.declarations || sym.declarations.length === 0) {
          // An undeclared symbol, e.g. "undefined".
          return;
        }
        const name = this.getSymbolName(sym, TSNamespace.VALUE);
        this.emitEdge(this.newAnchor(node), 'ref', name);
        return;
      default:
        // Use default recursive processing.
        return ts.forEachChild(node, n => this.visit(n));
    }
  }

  /** index is the main entry point, starting the recursive visit. */
  index() {
    this.emitFact(this.kFile, 'node/kind', 'file');
    this.emitFact(this.kFile, 'text', this.file.text);

    this.emitModuleAnchor(this.file);

    ts.forEachChild(this.file, n => this.visit(n));
  }
}

/**
 * index indexes a TypeScript program, producing Kythe JSON objects for the
 * source files in the specified paths.
 *
 * (A ts.Program is a configured collection of parsed source files, but
 * the caller must specify the source files within the program that they want
 * Kythe output for, because e.g. the standard library is contained within
 * the Program and we only want to process it once.)
 *
 * @param compilationUnit A VName for the entire compilation, containing e.g.
 *     corpus name.
 * @param pathVNames A map of file path to path-specific VName.
 * @param emit If provided, a function that receives objects as they are
 *     emitted; otherwise, they are printed to stdout.
 * @param readFile If provided, a function that reads a file as bytes to a
 *     Node Buffer.  It'd be nice to just reuse program.getSourceFile but
 *     unfortunately that returns a (Unicode) string and we need to get at
 *     each file's raw bytes for UTF-8<->UTF-16 conversions.
 */
export function index(
    vname: VName, pathVNames: Map<string, VName>, paths: string[],
    program: ts.Program, emit?: (obj: {}) => void,
    readFile: (path: string) => Buffer = fs.readFileSync) {
  // Note: we only call getPreEmitDiagnostics (which causes type checking to
  // happen) on the input paths as provided in paths.  This means we don't
  // e.g. type-check the standard library unless we were explicitly told to.
  const diags = new Set<ts.Diagnostic>();
  for (const path of paths) {
    for (const diag of ts.getPreEmitDiagnostics(
             program, program.getSourceFile(path))) {
      diags.add(diag);
    }
  }
  if (diags.size > 0) {
    const message = ts.formatDiagnostics(Array.from(diags), {
      getCurrentDirectory() {
        return program.getCompilerOptions().rootDir!;
      },
      getCanonicalFileName(fileName: string) {
        return fileName;
      },
      getNewLine() {
        return '\n';
      },
    });
    throw new Error(message);
  }

  const offsetTables = new Map<string, utf8.OffsetTable>();
  function getOffsetTable(path: string): utf8.OffsetTable {
    let table = offsetTables.get(path);
    if (!table) {
      const buf = readFile(path);
      table = new utf8.OffsetTable(buf);
      offsetTables.set(path, table);
    }
    return table;
  }

  for (const path of paths) {
    const sourceFile = program.getSourceFile(path);
    if (!sourceFile) {
      throw new Error(`requested indexing ${path} not found in program`);
    }
    const visitor =
        new Vistor(vname, pathVNames, program, sourceFile, getOffsetTable);
    if (emit != null) {
      visitor.emit = emit;
    }
    visitor.index();
  }
}

/**
 * loadTsConfig loads a tsconfig.json from a path, throwing on any errors
 * like "file not found" or parse errors.
 */
export function loadTsConfig(
    tsconfigPath: string, projectPath: string,
    host: ts.ParseConfigHost = ts.sys): ts.ParsedCommandLine {
  projectPath = path.resolve(projectPath);
  const {config: json, error} = ts.readConfigFile(tsconfigPath, host.readFile);
  if (error) {
    throw new Error(ts.formatDiagnostics([error], ts.createCompilerHost({})));
  }
  const config = ts.parseJsonConfigFileContent(json, host, projectPath);
  if (config.errors.length > 0) {
    throw new Error(
        ts.formatDiagnostics(config.errors, ts.createCompilerHost({})));
  }
  return config;
}

function main(argv: string[]) {
  if (argv.length < 1) {
    console.error('usage: indexer path/to/tsconfig.json [PATH...]');
    return 1;
  }

  const config = loadTsConfig(argv[0], path.dirname(argv[0]));
  let inPaths = argv.slice(1);
  if (inPaths.length === 0) {
    inPaths = config.fileNames;
  }

  // This program merely demonstrates the API, so use a fake corpus/root/etc.
  const compilationUnit: VName = {
    corpus: 'corpus',
    root: '',
    path: '',
    signature: '',
    language: '',
  };
  const program = ts.createProgram(inPaths, config.options);
  index(compilationUnit, new Map(), inPaths, program);
  return 0;
}

if (require.main === module) {
  // Note: do not use process.exit(), because that does not ensure that
  // process.stdout has been flushed(!).
  process.exitCode = main(process.argv.slice(2));
}
