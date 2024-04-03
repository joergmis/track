# clockodo

The [api documentation](https://www.clockodo.com/en/api/) is not very detailed
regarding how the different endpoints play together. Here is a short overview
based on the references in the documentation and the hints from the UI.

* There are different kind of services globally defined; choosing one is 
  required when adding an entry. It looks like they are **not** tied to a
  project or customer
* as expected; one customer has many projects
* `EntriesTexts` are not really used - there is the optional `text` property on
  `Entry` which is used in the UI (and therefore, it should probably be used)
