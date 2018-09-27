package pkg;

public class Selectors {
  //- @field defines/binding Field
  //- @String ref String
  String field;

  //- @Optional ref OptionalAbs
  //- @maybe defines/binding Param
  public String m(Optional<String> maybe) {
    //- @maybe ref Param
    //- @isPresent ref _IsPresentMethod
    if (maybe.isPresent()) {
      //- @maybe ref Param
      //- @get ref _GetMethod
      //- @field ref Field
      //- @this ref This
      this.field = maybe.get();
    }
    //- @this ref This
    //- @m2 ref M2Method
    //- @toString ref _ToStringMethod
    return this.m2().toString();
  }

  //- @m2 defines/binding M2Method
  private String m2() {
    //- @field ref Field
    //- @this ref This
    return this.field;
  }

  //- @String ref String
  //- @"java.lang" ref _JavaLangPackage
  java.lang.String m3() {
    return null;
  }

  //- @Optional defines/binding OptionalAbs
  //- Optional childof OptionalAbs
  //- Optional.node/kind interface
  private static interface Optional<T> {
    public T get();
    public boolean isPresent();
  }
}
