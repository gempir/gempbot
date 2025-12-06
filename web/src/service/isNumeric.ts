export function isNumeric(str: string): boolean {
  if (typeof str !== "string") return false; // we only process strings!
  // @ts-expect-error
  return (
    !Number.isNaN(str) && // use type coercion to parse the _entirety_ of the string (`parseFloat` alone does not do this)...
    !Number.isNaN(parseFloat(str))
  ); // ...and ensure strings of whitespace fail
}
