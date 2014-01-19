package testlib

import (
	"fmt"
	"image/png"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
)

func (t *TestSuite) TestDraw() {
	request := mandala.LoadAssetRequest{
		Filename: "res/drawable/expected.png",
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	if buffer != nil {
		exp, err := png.Decode(buffer)
		t.True(err == nil, "An error occured during png decoding")

		distance := imagetest.CompareDistance(exp, <-t.testDraw)
		t.True(distance < 0.1, fmt.Sprintf("Distance is %f", distance))
	}
}
