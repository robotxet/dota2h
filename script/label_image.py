#!/usr/bin/python
import sys
import tensorflow as tf

def calc_score(img_path, img_file):
    image_path = img_file
    # Read in the image_data
    image_data = tf.gfile.FastGFile(image_path, 'rb').read()

    # Loads label file, strips off carriage return
    label_lines = [line.rstrip() for line 
                       in tf.gfile.GFile(img_path + "/retrained_labels.txt")]

    # Unpersists graph from file
    with tf.gfile.FastGFile(img_path + "/retrained_graph.pb", 'rb') as f:
        graph_def = tf.GraphDef()
        graph_def.ParseFromString(f.read())
        _ = tf.import_graph_def(graph_def, name='')

    result = ""

    with tf.Session() as sess:
        # Feed the image_data as input to the graph and get first prediction
        softmax_tensor = sess.graph.get_tensor_by_name('final_result:0')

        predictions = sess.run(softmax_tensor, \
                 {'DecodeJpeg/contents:0': image_data})

        # Sort to show labels of first prediction in order of confidence
        top_k = predictions[0].argsort()[-len(predictions[0]):][::-1]

        for node_id in top_k:
            human_string = label_lines[node_id]
            score = predictions[0][node_id]
            result += '%s (score = %.5f)\n' % (human_string, score)
    return result

def main():
    print calc_score(sys.argv[1], sys.argv[2])

if __name__ == "__main__":
    main()