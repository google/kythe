// Tests the specification of VNames for symbol declarations

// SourceFile
//- FileModule=VName("module", _, _, "testdata/declaration_spec", "typescript").node/kind record
//- FileModuleAnchor.node/kind anchor
//- FileModuleAnchor./kythe/loc/start 0
//- FileModuleAnchor./kythe/loc/end 1
//- FileModuleAnchor defines/binding FileModule

// NamespaceImport
//- @NspI defines/binding VName("NspI", _, _, "testdata/declaration_spec", "typescript")
import * as NspI from './declaration';

// ExportAssignment
//- @default defines/binding VName("default", _, _, "testdata/declaration_spec", "typescript")
export default NspI;

// ClassDeclaration
//- @C defines/binding VName("C#type", _, _, "testdata/declaration_spec", "typescript")
class C {
  // PropertyDeclaration
  //- @property defines/binding VName("C.property", _, _, "testdata/declaration_spec", "typescript")
  property = 0;

  // MethodDeclaration
  //- @method defines/binding VName("C.method", _, _, "testdata/declaration_spec", "typescript")
  method() {}

  // Constructor
  //- @constructor defines/binding VName("C:ctor", _, _, "testdata/declaration_spec", "typescript")
  //- @pProp defines/binding VName("C.pProp", _, _, "testdata/declaration_spec", "typescript")
  constructor(private pProp: number) {
    //- @a defines/binding VName("C.constructor.a", _, _, "testdata/declaration_spec", "typescript")
    let a;
  }

  // GetAccessor
  //- @prop defines/binding VName("C.prop", _, _, "testdata/declaration_spec", "typescript")
  get prop() {
    return this.property;
  }

  // SetAccessor
  //- @prop defines/binding VName("C.prop:setter", _, _, "testdata/declaration_spec", "typescript")
  set prop(nProp) {
    this.property = nProp;
  }
}

// ClassDeclaration with no constructor
//- @CC defines/binding VName("CC:ctor", _, _, "testdata/declaration_spec", "typescript")
class CC {}

// EnumDeclaration
//- @E defines/binding VName("E", _, _, "testdata/declaration_spec", "typescript")
//- @E defines/binding VName("E#type", _, _, "testdata/declaration_spec", "typescript")
enum E {
  // EnumMember
  //- @EnumMember defines/binding VName("E.EnumMember", _, _, "testdata/declaration_spec", "typescript")
  EnumMember = 0
}

// FunctionDeclaration
//- @#1"fun" defines/binding VName("fun", _, _, "testdata/declaration_spec", "typescript")
function fun(
    // Parameter
    //- @param defines/binding VName("fun.param", _, _, "testdata/declaration_spec", "typescript")
    param: number) {}

// InterfaceDeclaration
//- @B defines/binding VName("B#type", _, _, "testdata/declaration_spec", "typescript")
interface B {
  // PropertySignature
  //- @pSig defines/binding VName("B.pSig", _, _, "testdata/declaration_spec", "typescript")
  pSig: number;

  // MethodSignature
  //- @mSig defines/binding VName("B.mSig", _, _, "testdata/declaration_spec", "typescript")
  mSig(): void;
}

// VariableDeclaration
//- @v defines/binding VName("v", _, _, "testdata/declaration_spec", "typescript")
let v = {
  // PropertyAssignment
  // TODO: the signature here should be something like `block0.prop`, but
  // anonymous block names are not well-defined by the spec yet.
  //- @prop defines/binding VName(_, _, _, "testdata/declaration_spec", "typescript")
  prop: 0
};

// TypeAliasDeclaration
//- @AliasArray defines/binding VName("AliasArray#type", _, _, "testdata/declaration_spec", "typescript")
type AliasArray<
    // TypeParameter
    //- @#0"T" defines/binding VName("AliasArray.T#type", _, _, "testdata/declaration_spec", "typescript")
    T> = Array<T>;

//- @arrowFun defines/binding VName("arrowFun", _, _, "testdata/declaration_spec", "typescript")
const arrowFun = () => {
  // Arrow function scope name is not well-defined.
  //- @anonArrowFunDecl defines/binding VName(_, _, _, "testdata/declaration_spec", "typescript")
  let anonArrowFunDecl;
};

{
  // Anonymous block scope name is not well-defined.
  //- @anonBlockDecl defines/binding VName(_, _, _, "testdata/declaration_spec", "typescript")
  let anonBlockDecl;
}
