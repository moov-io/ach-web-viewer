<table style="border-spacing: 10px; margin-bottom: 10px; min-width:100%;">
  <!-- FileHeader -->
  <caption id="file-header">File Header</caption>
  <tr>
    <td>Origin</td>
    <td>OriginName</td>
    <td>Destination</td>
    <td>DestinationName</td>
    <td>FileCreationDate</td>
    <td>FileCreationTime</td>
  </tr>
  <tr>
    <td>{{ .Header.ImmediateOrigin }} </td>
    <td>{{ .Header.ImmediateOriginName }} </td>
    <td>{{ .Header.ImmediateDestination }} </td>
    <td>{{ .Header.ImmediateDestinationName }} </td>
    <td>{{ .Header.FileCreationDate }} </td>
    <td>{{ .Header.FileCreationTime }} </td>
  </tr>
</table>

<table style="border-spacing: 10px; margin-bottom: 10px; min-width:100%;">
  <!-- Batches -->
  <caption id="batches">Batches</caption>
  {{ range $batch := .Batches }}
  <tr>
    <td>BatchNumber</td>
    <td>SECCode</td>
    <td>ServiceClassCode</td>
    <td><a href="?format=table&sort=CompanyName&order=asc">CompanyName</a></td>
    <td>DiscretionaryData</td>
    <td><a href="?format=table&sort=CompanyIdentification&order=asc">Identification</a></td>
    <td>EntryDescription</td>
    <td>EffectiveEntryDate</td>
    <td>DescriptiveDate</td>
  </tr>
  <tr>
    <td>{{ $batch.Header.BatchNumber }}</td>
    <td>{{ $batch.Header.StandardEntryClassCode }}</td>
    <td>{{ $batch.Header.ServiceClassCode }}</td>
    <td>{{ $batch.Header.CompanyName }}</td>
    <td>{{ $batch.Header.CompanyDiscretionaryData }}</td>
    <td>{{ $batch.Header.CompanyIdentification }}</td>
    <td>{{ $batch.Header.CompanyEntryDescription }}</td>
    <td>{{ $batch.Header.EffectiveEntryDate }}</td>
    <td>{{ $batch.Header.CompanyDescriptiveDate }}</td>
  </tr>
  <tr>
    <td style="padding-left:20px;">TransactionCode</td>
    <td>RDFIIdentification</td>
    <td>AccountNumber</td>
    <td><a href="?format=table&sort=Amount&order=desc">Amount</a></td>
    <td><a href="?format=table&sort=IndividualName&order=asc">Name</a></td>
    <td><a href="?format=table&sort=TraceNumber&order=asc">TraceNumber</a></td>
    <td>Category</td>
  </tr>
  {{ range $entry := $batch.GetEntries }}
  <tr>
    <td style="padding-left:20px;">{{ $entry.TransactionCode }}</td> <!-- TODO: view helper -->
    <td>{{ $entry.RDFIIdentification }}</td>
    <td>{{ $entry.DFIAccountNumber }}</td> <!-- TODO(adam): masking -->
    <td>{{ $entry.Amount }}</td>           <!-- TODO(adam): human readable formatting -->
    <td>{{ $entry.IndividualName }}</td>
    <td>{{ $entry.TraceNumber }}</td>
    <td>{{ $entry.Category }}</td>
  </tr>
  <!-- TODO(adam): print addenda records and intent -->
  {{ end }}
  {{ end }}
</table>

<table style="border-spacing: 10px; margin-bottom: 10px; min-width:100%;">
  <!-- FileControl -->
  <caption id="file-control">File Control</caption>
  <tr>
    <td>BatchCount</td>
    <td>BlockCount</td>
    <td>EntryAddendaCount</td>
    <td>TotalDebitAmount</td>
    <td>TotalCreditAmount</td>
  </tr>
  <tr>
    <td>{{ .Control.BatchCount }}</td>
    <td>{{ .Control.BlockCount }}</td>
    <td>{{ .Control.EntryAddendaCount }}</td>
    <td>{{ .Control.TotalDebitEntryDollarAmountInFile }}</td>
    <td>{{ .Control.TotalCreditEntryDollarAmountInFile }}</td>
  </tr>
</table>
