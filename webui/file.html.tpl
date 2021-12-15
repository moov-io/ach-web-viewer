<!doctype html>
<html>
  <head>
    <title>{{ .Filename }} | ACH viewer</title>
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
      <a href="{{ .BaseURL }}"><- Back</a>
      <h1>{{ .Filename }}</h1>
    </header>
    {{ if ne .Valid nil }}
      <div class="error">
        {{ .Valid }}
      </div>
    {{ end }}
    <main>
      <pre>{{ .Contents }}</pre>
      <hr />
      <table>
        <thead>
          <tr colspan="2">
            <td id="metadata-header">File Metadata</td>
          </tr>
        </thead>
        {{ range $key, $value := .Metadata }}
        <tr>
          <td class="metadata-key">{{ $key }}</td>
          <td>{{ $value }}</td>
        </tr>
        {{ end }}
      </table>
    </main>
  </body>
</html>
