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

	st-proc [-v] [<hostname> [<port>]]

## Description

``st-proc`` connects to a server at ``<hostname>:<port>`` that supplies a data
message requiring processing, decodes the data and writes to stdout.  It
blocks waiting for data from the server, stop with ``Ctrl-C`` or ``kill(1)``.

* ``<hostname>``
  ``hostname`` can be a numerical IP or symbolic hostname.  Default is
  ``data.salad.com``.

* ``<port>``
  ``port`` is a positive integer.  Default is ``5000``.

### flight Package

All of the specifics of handling the flight data are in the ``flight`` package.
This primarily consists of the decoder and a header object that includes
JSON ``Marshal()`` and ``Unmarshal()`` support for base64 encoding.

* Note that the specifications say the size of the floats are 4 bytes but also describes
  them as IEEE-754 64 bit (binary64). The binary buffer included uses 8 bytes for the
  float values so that is what is used in the code.

### cmd Package

The ``cmd`` package contains the main processing loop and router connection handling.

The processing loop remains single-threaded since there is no requirement for parallel processing.

The ``RouterConn`` struct encapsulates the network info and configuration for
connecting to the router and the retry logic for the ``Connect()`` method.  It uses
``io.Copy`` and ``bytes.Buffer`` to dynamically size the receive buffer since
the simple protocol only sends a single stream of data then closes the connection.

## Building

* Get the source:

    ```bash
	git clone https://github.com/dtroyer/st-proc
	cd st-proc
	make setup
    ```

* Build, the executable will be in bin/st-proc:

    ```bash
	make build
    ```

* Run unit tests:

    ```bash
    make test
    ```

* Run locally: We can use ``netcat`` and ``base64`` to simulate the server and
  run the processor locally for testing.  Two shell sessions are required.  To
  start the server, run in one session:

    ```bash
    base64 -d files/testPacket1 | nc -l 5000
    ```

  And in the other session run the processor client:

    ```bash
    bin/st-proc localhost
    ```

  Compare the JSON output with the testJson1 string in ``flight/message_test.go``
  or with the string in the original spec.  Note that the ``header`` field is
  included in this implementation.

  All 3 of the testPacketN messages are included in ``files/`` base64 encoded.

## See Also

Source [on Github](https://github.com/dtroyer/st-proc)

## Author

Dean Troyer <dt-github@xr7.org>

## Known Issues

* The JSON output produces the ``header`` value as an integer array rather than a
  base64-encoded string.  ``message.UnmarshalJSON()`` is not being called as expected.

    [![Go Report Card](https://goreportcard.com/badge/github.com/dtroyer/st-proc?style=flat-square)](https://goreportcard.com/report/github.com/dtroyer/st-proc)
