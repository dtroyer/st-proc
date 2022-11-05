# Message Processor

## Requirements

* Connect to an upstream router and wait for messages
  * Retry if connection closed with no message received
* Receive a message
* Decode message into a struct
* Write struct to stdout
* Repeat
* Strings are encoded as UTF-8
* Cleanly handle expected signals

## Usage

	st-proc [<hostname> [<port>]]

## Description

``st-proc`` connects to a server at ``<hostname>:<port>`` that supplies a data
message requiring processing, decodes the data and writes to stdout.  It
blocks waiting for data from the server, stop with ``Ctrl-C`` or ``kill(1)``.

* ``<hostname>``
  ``hostname`` can be a numerical IP or symbolic hostname.  Default is
  ``data.salad.com``.

* ``<port>``
  ``port`` is a positive integer.  Default is ``5000``.

## See Also

Source [on Github](https://github.com/dtroyer/st-proc)

## Author

Dean Troyer <dt@xr7.org>
