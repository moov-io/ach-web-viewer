<!doctype html>
<html>
  <head>
    <title>ACH file viewer</title>
    <link rel="stylesheet" href="{{ .BaseURL }}style.css">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Manrope:wght@400;500;700&display=swap" rel="stylesheet">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,viewport-fit=cover">
    <meta name="theme-color" media="(prefers-color-scheme: light)" content="#fefefe">
    <meta name="theme-color" media="(prefers-color-scheme: dark)" content="#202020">
  </head>
  <body>
    <header>
      <h1>ACH file viewer</h1>
    </header>
    <main class="list">
      {{ range $source := .Sources }}
        <div class="source"><strong>{{ $source.ID }}</strong> ({{ $source.Type }})</div>
        {{ range $group := $source.Groups }}
          <div class="date">{{ dateTime $group.DateTime "January 2, 2006" }}</div>
          {{ range $file := $group.Files }}
            <a href="{{ $file.Path }}">
              <svg class="icon" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m18 2h-12a2 2 0 0 0 -2 2v16a2 2 0 0 0 2 2h6.76a3 3 0 0 0 2.12-.88l4.24-4.24a3 3 0 0 0 .88-2.12v-10.76a2 2 0 0 0 -2-2zm0 12h-5a1 1 0 0 0 -1 1v5h-6v-16h12zm-9.5-6h7a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0 -.5-.5h-7a.5.5 0 0 0 -.5.5v1a.5.5 0 0 0 .5.5zm0 4h4a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0 -.5-.5h-4a.5.5 0 0 0 -.5.5v1a.5.5 0 0 0 .5.5z"></path></svg>
              <div>
                <span class="clean-path">{{ $file.CleanPath }}</span>
                {{ $file.Filename }}
              </div>
            </a>
          {{ end }}
        {{ end }}
      {{ end }}

      <span>
        <a href="{{ startDateParam .Options.TimeRangeMin }}">Previous</a> / <a href="{{ endDateParam .Options.TimeRangeMax }}">Next</a>
      </span>
    </main>
  </body>
</html>
