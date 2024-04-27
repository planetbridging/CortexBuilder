# NEAT (NeuroEvolution of Augmenting Topologies)

NEAT is a sophisticated type of TWEANN (Topology and Weight Evolving Artificial Neural Network) that utilizes genetic algorithms to evolve artificial neural networks. Developed by Kenneth O. Stanley, this innovative framework is designed to optimize both the structure (topology) and the synaptic strengths (weights) of neural networks simultaneously. By dynamically adapting network architectures and connection weights, NEAT efficiently tackles complex problems by evolving networks that are intricately tailored to specific tasks. This dual evolution approach allows NEAT to discover novel network topologies along with optimal weight configurations, making it a powerful tool for developing advanced neural models that can generalize well across varied datasets.

## Key Components and Concepts

### 1. **Genome Representation**

- **Nodes**: Represents neurons in the neural network, categorized as input, hidden, or output.
- **Connections**: Represents synapses between neurons, each with a weight and an enabled/disabled state.
- **Genes**: Each connection gene specifies the in-node, out-node, weight, its enabled state, and an innovation number for tracking historical changes.

### 2. **Population Initialization**

- Begins with minimal structures (just input and output nodes) to facilitate the evolution of complex structures over time.

### 3. **Fitness Evaluation**

- **Training and Validation Split**: Implements holdout validation where the dataset is split into training (70%) and validation (30%) segments.
- **Performance Measurement**: Each genome is evaluated on its ability to perform a given task using the training data, but its fitness is primarily determined by its performance on the validation data to ensure generalization.

### 4. **Selection**

- Uses methods like tournament selection or roulette wheel selection based on fitness to choose genomes for reproduction.

### 5. **Crossover (Recombination)**

- Combines genomes from two parents, respecting the historical origins of genes to maintain structural integrity.

### 6. **Mutation**

- **Weights Mutation**: Modifies the weights of connections to adapt and refine the network’s responses to inputs.
- **Bias Mutation**: Adjusts the biases of neurons to fine-tune the activation potential, enhancing the network's ability to fit complex patterns.
- **Add Node Mutation**: Inserts a new node by splitting an existing connection, increasing the network’s depth and potential for complexity.
- **Add Connection Mutation**: Creates a new connection between previously unconnected nodes, expanding the network’s capacity for diverse interactions.
- **Connection Enable/Disable**: Toggles the enabled state of connections, allowing the network to experiment with different neural pathways without permanent structural changes.(Will be developed later)
- **Add Layer Mutation**: Introduces entirely new layers to the network, significantly enhancing its depth and functional complexity.
- **Activation Function Mutation**: Alters the activation function of nodes to better suit different types of data processing needs, adapting to the specific characteristics of the input data.
- **Node Type Mutation**: Switches node types (e.g., from standard neurons to LSTM units or attention mechanisms), enabling the network to handle temporal dynamics or focus on relevant input features effectively.
  (Will be developed later)

Will be adding the disabled node and node type later on because of the complexity and need to get to a working status first.

### 7. **Speciation**

- Groups similar topologies into species to protect innovation and ensure that new structures are not immediately outcompeted.

### 8. **Reproduction**

- Selects the fittest individuals within each species for reproduction, replacing less fit individuals in the population, and encourages diversity through penalties on stagnant species.

### 9. **Parameter Settings**

- Includes mutation rates, population size, tournament size, and speciation thresholds, all needing careful adjustment based on specific problems.

### 10. **Termination**

- Ends based on satisfaction criteria like achieving a certain fitness level or after a predetermined number of generations.

## Practical Implementation Tips

- **Logging and Analysis**: Essential for tracking changes and understanding evolutionary dynamics.
- **Parallelization**: Speeds up fitness evaluations across the population.
- **Dynamic Parameters**: Adjust mutation rates and other parameters in response to evolutionary progress to maintain diversity and drive convergence.

## Conclusion

Building a NEAT implementation requires a deep understanding of both evolutionary algorithms and neural network principles. With tools such as NEAT-Python, developers can create robust solutions tailored to complex problem-solving scenarios. This framework emphasizes the importance of generalization through holdout validation, ensuring that evolved networks perform well on unseen data.

**What a feedforward neural network with multiple activation functions can do:**
(This is what is currently aiming for then working on more advanced)

1. Image Classification
2. Object Detection
3. Image Segmentation
4. Zero-Shot Image Classification
5. Text Classification
6. Token Classification
7. Table Question Answering
8. Question Answering
9. Zero-Shot Classification
10. Feature Extraction (Image and Text)
11. Tabular Classification
12. Tabular Regression

**What a feedforward neural network with multiple activation functions cannot do:**

1. Image-Text-to-Text
2. Visual Question Answering
3. Document Question Answering
4. Depth Estimation
5. Text-to-Image
6. Image-to-Text
7. Image-to-Image
8. Image-to-Video
9. Unconditional Image Generation
10. Video Classification
11. Text-to-Video
12. Mask Generation
13. Zero-Shot Object Detection
14. Text-to-3D
15. Image-to-3D
16. Translation
17. Summarization
18. Text Generation
19. Text2Text Generation
20. Fill-Mask
21. Sentence Similarity
22. Text-to-Speech
23. Text-to-Audio
24. Automatic Speech Recognition
25. Audio-to-Audio
26. Audio Classification
27. Voice Activity Detection
28. Reinforcement Learning
29. Robotics
30. Graph Machine Learning

**What a feedforward neural network with self-attention can do:**

1. Image Classification
2. Object Detection
3. Image Segmentation
4. Zero-Shot Image Classification
5. Text Classification
6. Token Classification
7. Table Question Answering
8. Question Answering
9. Zero-Shot Classification
10. Feature Extraction (Image and Text)
11. Tabular Classification
12. Tabular Regression
13. Translation
14. Summarization
15. Text Generation
16. Text2Text Generation
17. Fill-Mask
18. Sentence Similarity

**What a feedforward neural network with self-attention cannot do:**

1. Image-Text-to-Text
2. Visual Question Answering
3. Document Question Answering
4. Depth Estimation
5. Text-to-Image
6. Image-to-Text
7. Image-to-Image
8. Image-to-Video
9. Unconditional Image Generation
10. Video Classification
11. Text-to-Video
12. Mask Generation
13. Zero-Shot Object Detection
14. Text-to-3D
15. Image-to-3D
16. Text-to-Speech
17. Text-to-Audio
18. Automatic Speech Recognition
19. Audio-to-Audio
20. Audio Classification
21. Voice Activity Detection
22. Reinforcement Learning
23. Robotics
24. Graph Machine Learning
