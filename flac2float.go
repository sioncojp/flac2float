package flac2float

import (
	"io"
	"math"

	"azul3d.org/engine/audio"

	_ "azul3d.org/engine/audio/flac"
	_ "azul3d.org/engine/audio/wav"
)

var (
	ErrFormat        = audio.ErrFormat
	ErrInvalidData   = audio.ErrInvalidData
	ErrUnexpectedEOS = audio.ErrUnexpectedEOS
)

type SampleReduceFunc func(samples audio.Float64) float64

type Sound struct {
	r          io.Reader
	resolution uint
	sampleFn   SampleReduceFunc
}

func RMSF64Samples(samples audio.Float64) float64 {
	var sumSquare float64
	for i := range samples {
		sumSquare += math.Pow(samples.At(i), 2)
	}

	return math.Sqrt(sumSquare / float64(samples.Len()))
}

func New(r io.Reader, i uint) *Sound {
	s := &Sound{
		r:          r,
		resolution: i,
		sampleFn:   RMSF64Samples,
	}
	return s
}

func (s *Sound) ReadSound() ([]float64, error) {
	decoder, _, err := audio.NewDecoder(s.r)
	if err != nil {
		// Unknown format
		if err == audio.ErrFormat {
			return nil, ErrFormat
		}

		// Invalid data
		if err == audio.ErrInvalidData {
			return nil, ErrInvalidData
		}

		// Unexpected end-of-stream
		if err == audio.ErrUnexpectedEOS {
			return nil, ErrUnexpectedEOS
		}

		// All other errors
		return nil, err
	}
	var computed []float64
	var value float64
	config := decoder.Config()
	samples := make(audio.Float64, uint(config.SampleRate*config.Channels)/s.resolution)

	for {
		_, err := decoder.Read(samples)
		if err != nil && err != audio.EOS {
			return nil, err
		}
		value = s.sampleFn(samples)
		computed = append(computed, value)
		if err == audio.EOS {
			break
		}
	}

	// Return slice of computed values
	return computed, nil
}
