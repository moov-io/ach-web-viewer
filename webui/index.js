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
};

function filterFileListing(pattern) {
  const url = new URL(window.location.href);
  url.searchParams.set("pattern", pattern);

  window.location.href = url.toString(); // redirect
}

function populateCurrentPatternForm() {
  const url = new URL(window.location.href);

  const currentPattern = url.searchParams.get("pattern");
  if (currentPattern) {
    var filenameElm = document.querySelector("#filename");
    filenameElm.value = currentPattern;
  }
}
