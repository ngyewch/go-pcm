package pcm

import "time"

type SourceConverter struct {
	source               Source
	encoding             Encoding
	bytesPerSourceSample int
	bytesPerTargetSample int
}

func NewSourceConverter(source Source, encoding Encoding) *SourceConverter {
	return &SourceConverter{
		source:               source,
		encoding:             encoding,
		bytesPerSourceSample: source.Encoding().BytesPerSample(),
		bytesPerTargetSample: encoding.BytesPerSample(),
	}
}

func (reader *SourceConverter) Encoding() Encoding {
	return reader.encoding
}

func (reader *SourceConverter) NumChannels() uint16 {
	return reader.source.NumChannels()
}

func (reader *SourceConverter) SampleRate() uint32 {
	return reader.source.SampleRate()
}

func (reader *SourceConverter) Duration() time.Duration {
	return reader.source.Duration()
}

func (reader *SourceConverter) Read(data []byte) (int, error) {
	if reader.source.Encoding() == reader.encoding {
		return reader.source.Read(data)
	}
	targetSamples := len(data) / reader.bytesPerTargetSample
	sourceData := make([]byte, targetSamples*reader.bytesPerSourceSample)
	sourceLen, err := reader.source.Read(sourceData)
	if err != nil {
		return 0, err
	}
	sourceSampleCount := sourceLen / reader.bytesPerSourceSample
	targetData := data[:sourceSampleCount*reader.bytesPerTargetSample]
	err = reader.source.Encoding().ConvertSlice(sourceData[:sourceLen], reader.encoding, targetData)
	if err != nil {
		return 0, err
	}
	return len(targetData), nil
}
