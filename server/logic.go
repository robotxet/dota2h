package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
    "github.com/tensorflow/tensorflow/tensorflow/go/op"
)

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
func constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {
    // Some constants specific to the pre-trained model at:
    // https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip
    //
    // - The model was trained after with images scaled to 224x224 pixels.
    // - The colors, represented as R, G, B in 1-byte each were converted to
    //   float using (value - Mean)/Scale.
    const (
        H, W  = 224, 224
        Mean  = float32(117)
        Scale = float32(1)
    )
    // - input is a String-Tensor, where the string the JPEG-encoded image.
    // - The inception model takes a 4D tensor of shape
    //   [BatchSize, Height, Width, Colors=3], where each pixel is
    //   represented as a triplet of floats
    // - Apply normalization on each pixel and use ExpandDims to make
    //   this single image be a "batch" of size 1 for ResizeBilinear.
    s := op.NewScope()
    input = op.Placeholder(s, tf.String)
    output = op.Div(s,
        op.Sub(s,
            op.ResizeBilinear(s,
                op.ExpandDims(s,
                    op.Cast(s,
                        op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)), tf.Float),
                    op.Const(s.SubScope("make_batch"), int32(0))),
                op.Const(s.SubScope("size"), []int32{H, W})),
            op.Const(s.SubScope("mean"), Mean)),
        op.Const(s.SubScope("scale"), Scale))
    graph, err = s.Finalize()
    return graph, input, output, err
}

func modelFiles(dir string) (modelfile, labelsfile string, err error) {
    var (
        model   = filepath.Join(dir, "retrained_graph.pb")
        labels  = filepath.Join(dir, "retrained_labels.txt")
    )
    if filesExist(model, labels) == nil {
        return model, labels, nil
    }
    return model, labels, filesExist(model, labels)
}

func filesExist(files ...string) error {
    for _, f := range files {
        if _, err := os.Stat(f); err != nil {
            return fmt.Errorf("unable to stat %s: %v", f, err)
        }
    }
    return nil
}

func printBestLabel(probabilities []float32, labelsFile string) {
    bestIdx := 0
    for i, p := range probabilities {
        if p > probabilities[bestIdx] {
            bestIdx = i
        }
    }
    // Found the best match. Read the string from labelsFile, which
    // contains one line per label.
    file, err := os.Open(labelsFile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    var labels []string
    for scanner.Scan() {
        labels = append(labels, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Printf("ERROR: failed to read %s: %v", labelsFile, err)
    }
    log.Printf("BEST MATCH: (%2.0f%% likely) %s\n", probabilities[bestIdx]*100.0, labels[bestIdx])
}

// Convert the image in filename to a Tensor suitable as input to the Inception model.
func makeTensorFromImage(filename string) (*tf.Tensor, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    // DecodeJpeg uses a scalar String-valued tensor as input.
    tensor, err := tf.NewTensor(string(bytes))
    if err != nil {
        return nil, err
    }
    // Construct a graph to normalize the image
    graph, input, output, err := constructGraphToNormalizeImage()
    if err != nil {
        return nil, err
    }
    // Execute that graph to normalize this one image
    session, err := tf.NewSession(graph, nil)
    if err != nil {
        return nil, err
    }
    defer session.Close()
    normalized, err := session.Run(
        map[tf.Output]*tf.Tensor{input: tensor},
        []tf.Output{output},
        nil)
    if err != nil {
        return nil, err
    }
    return normalized[0], nil
}


func getPredictions(imagePath string, tensorflowPath string) string {
	modelfile, labelfile, err := modelFiles(tensorflowPath)
	if err != nil {
		log.Fatal(err)
	}
	model, err := ioutil.ReadFile(modelfile)
	if err != nil {
		log.Fatal(err)
	}
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Println("Failed to load graph")
		log.Fatal(err)
	}
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	tensor, err := makeTensorFromImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	probabilities := output[0].Value().([][]float32)[0]
	printBestLabel(probabilities, labelfile)
	return ""
}



// #!/usr/bin/python
// import sys
// import tensorflow as tf

// def calc_score(img_path, img_file):
//     image_path = img_file
//     # Read in the image_data
//     image_data = tf.gfile.FastGFile(image_path, 'rb').read()

//     # Loads label file, strips off carriage return
//     label_lines = [line.rstrip() for line 
//                        in tf.gfile.GFile(img_path + "/retrained_labels.txt")]

//     # Unpersists graph from file
//     with tf.gfile.FastGFile(img_path + "/retrained_graph.pb", 'rb') as f:
//         graph_def = tf.GraphDef()
//         graph_def.ParseFromString(f.read())
//         _ = tf.import_graph_def(graph_def, name='')

//     result = ""

//     with tf.Session() as sess:
//         # Feed the image_data as input to the graph and get first prediction
//         softmax_tensor = sess.graph.get_tensor_by_name('final_result:0')

//         predictions = sess.run(softmax_tensor, \
//                  {'DecodeJpeg/contents:0': image_data})

//         # Sort to show labels of first prediction in order of confidence
//         top_k = predictions[0].argsort()[-len(predictions[0]):][::-1]

//         for node_id in top_k:
//             human_string = label_lines[node_id]
//             score = predictions[0][node_id]
//             result += '%s:%.5f\n' % (human_string, score)
//     return result

// def main():
//     print calc_score(sys.argv[1], sys.argv[2])

// if __name__ == "__main__":
//     main()