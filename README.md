# Overview

This is a simple commandline utility app to create test bags for APTrust using
 the BagIt Specifications.

It will bag Fedora Datastream files or DSpace AIPs as nessisary.

See the APTrust BagIt profile document on Google Docs for moreinformation about
 the structure of APTrust bags.

This is not intended to be a community too, it is only to help create test bags
for trial ingests into APTrust

## DSpace AIPs

AIPs originating from DSpace will be uncompressed and stored in the bag data
directory to allow for the examination of the contents after processing.

## Fedora Objects

Datastreams from Fedora Objects are written as individual files into the
staging area.  Bag data directories will cointain all the individual 
datastream files for a Fedora Object, including any versioned datastream
files.