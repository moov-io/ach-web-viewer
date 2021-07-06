<!doctype html>
<html>
  <head>
    <title>{{ .Filename }} | ACH Viewer</title>
  </head>
  <body>
    <a href="{{ .BackURL }}">Back</a>
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
