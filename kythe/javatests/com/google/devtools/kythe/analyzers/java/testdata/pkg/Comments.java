package pkg;

//- @+11"java.lang.Integer" ref/doc IntegerClass
//- @+10String ref/doc StringClass
//- @+10Inner ref/doc InnerClass

//- DocComment.node/kind anchor
//- DocComment documents CommentsClass
//- DocComment.loc/start @^+4"/**"
//- DocComment.loc/end @$+6"*/"
//- @+6Comments defines/binding CommentsClass

/**
 * This is a Javadoc comment with links to {@link String}, {@link java.lang.Integer}, and
 * {@link Inner}.
 */
public class Comments {

  //- @+3"// inline comment here" documents FieldOne
  //- @+2fieldOne defines/binding FieldOne

  private static int fieldOne; // inline comment here

  //- @+3"// fieldTwo represents the universe" documents FieldTwo
  //- @+3fieldTwo defines/binding FieldTwo

  // fieldTwo represents the universe
  private static String fieldTwo;

  //- @+3"// This comments the Inner class." documents InnerClass
  //- @+3Inner defines/binding InnerClass

  // This comments the Inner class.
  public static class Inner {}

  //- @+5"// a second, weirdly-placed comment" documents InnerI
  //- @+5"// this also comments the interface" documents InnerI
  //- @+3"/* This comments the InnerI interface. */" documents InnerI
  //- @+3InnerI defines/binding InnerI

  /* This comments the InnerI interface. */ // a second, weirdly-placed comment
  static interface InnerI {} // this also comments the interface

  //- @+5InnerE ref/doc InnerE

  //- @+3"/** This comments the {@link InnerE} enum. */" documents InnerE
  //- @+3InnerE defines/binding InnerE

  /** This comments the {@link InnerE} enum. */
  static enum InnerE {}
  // TODO(schroederc): link JavaDoc
}
