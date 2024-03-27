// This script is injected into the page and is used to extract text from the page

() => {
  function textNodesUnder(el) {
    var n,
      a = [],
      walk = document.createTreeWalker(el, NodeFilter.SHOW_TEXT, null, false);
    while ((n = walk.nextNode())) a.push(n);
    return a;
  }

  return [
    ...new Set(
      textNodesUnder(document.body)
        .filter(
          (element) =>
            element.parentElement.tagName !== "SCRIPT" &&
            element.parentElement.tagName !== "STYLE" &&
            element.parentElement.tagName !== "NOSCRIPT" &&
            element.parentElement.tagName !== "OPTION",
        )
        .map((v) => v.nodeValue)
        .map((v) => v.trim())
        .filter((v) => !!v),
    ),
  ]
    .map((v) => v.substring(0, 400))
    .join(" ")
    .replaceAll("\n", "");
};
