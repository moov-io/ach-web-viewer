<!doctype html>
<html>
  <head>
    <title>{{ .Filename }} | ACH Viewer</title>
    <link rel="stylesheet" href="{{ .BaseURL }}style.css">
  </head>
  <body>
    <a href="{{ .BaseURL }}">Back</a>
    <br />
    <pre>{{ .Contents }}</pre>
    <br />
    {{ if eq .Valid nil }}
    Valid: true
    {{ else }}
    <strong>Validation Error</strong>:
    <br /><pre>{{ .Valid }}</pre>
    {{ end }}
  </body>
</html>
