package packets

import "io"

func readAll(r io.Reader, readers ...io.ReaderFrom) (int64, error) {
	var totalRead int64
	for _, reader := range readers {
		n, err := reader.ReadFrom(r)
		totalRead += n
		if err != nil {
			return totalRead, err
		}
	}
	return totalRead, nil
}

func writeAll(w io.Writer, writers ...io.WriterTo) (int64, error) {
	var totalWritten int64
	for _, writer := range writers {
		n, err := writer.WriteTo(w)
		totalWritten += n
		if err != nil {
			return totalWritten, err
		}
	}
	return totalWritten, nil
}
