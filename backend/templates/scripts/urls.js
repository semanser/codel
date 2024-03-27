// This script is injected into the page and is used to extract all urls from the page

() => {
  function extractUrlsFromLinks(el) {
    var links = Array.from(el.getElementsByTagName("a"));
    return links.map((link) => `[${link.textContent}](${link.href})`);
  }

  function extractUrlsFromAttributes(el) {
    var attributes = ["src", "href"];
    return attributes.map((attr) => {
      var elements = Array.from(el.querySelectorAll(`[${attr}]`));
      return elements.map((element) => `[${element.getAttribute(attr)}](${element.getAttribute(attr)})`);
    }).flat();
  }

  var linksArray = extractUrlsFromLinks(document);
  var attributesArray = extractUrlsFromAttributes(document);
  
  var allUrls = [...new Set([...linksArray, ...attributesArray])];

  return allUrls.join("\n ");
};
