---
'grafana-google-sheets-datasource': patch
---

Fix panic when a number-formatted cell has no computed value (e.g. a formula error like #DIV/0!, #N/A, #REF!)
