window.onload = function() {
  // Register a handler on the filter form
  var form = document.querySelector("#filter");
  form.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      event.preventDefault();

      var filenameElm = document.querySelector("#filename");
      filterFileListing(filenameElm.value);
    }
  });

  // Populate the filter with any current pattern from the URL
  populateCurrentPatternForm();

  // Update pagination links with the pattern
  updatePaginationLinks();
};

function filterFileListing(pattern) {
  const url = new URL(window.location.href);
  url.searchParams.set("pattern", pattern);

  window.location.href = url.toString(); // redirect
}

function getCurrentPattern() {
  const url = new URL(window.location.href);

  return url.searchParams.get("pattern");
}

function populateCurrentPatternForm() {
  const currentPattern = getCurrentPattern();
  if (currentPattern) {
    var filenameElm = document.querySelector("#filename");
    filenameElm.value = currentPattern;
  }
}

function updatePaginationLinks() {
  const currentPattern = getCurrentPattern();
  if (currentPattern) {
    var links = document.querySelectorAll("nav.pagination a");
    links.forEach(link => {
      link.href += "&pattern=" + currentPattern;
    });
  }
}
