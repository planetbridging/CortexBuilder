import React, { useState } from "react";
import {
  Button,
  FormControl,
  FormLabel,
  Input,
  NumberInput,
  NumberInputField,
  Radio,
  RadioGroup,
  Stack,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";

function OTraining() {
  const [networkType, setNetworkType] = useState("feedforward");
  const [spawnCount, setSpawnCount] = useState(100);
  const [additionalParam1, setAdditionalParam1] = useState("");
  const [additionalParam2, setAdditionalParam2] = useState("");
  const [additionalParam3, setAdditionalParam3] = useState("");
  const toast = useToast();

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await axios.post("http://localhost:4123/initialize", {
        networkType,
        spawnCount,
        additionalParam1,
        additionalParam2,
        additionalParam3,
      });
      console.log(response.data); // Handle the response as needed
      toast({
        title: "Training initialized",
        description:
          "The neural network training initialization process has started.",
        status: "success",
        duration: 5000,
        isClosable: true,
      });
    } catch (error) {
      console.error("Error:", error);
      toast({
        title: "Error",
        description: "Failed to initialize training.",
        status: "error",
        duration: 5000,
        isClosable: true,
      });
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <FormControl>
        <FormLabel>Network Type</FormLabel>
        <RadioGroup
          onChange={(e) => setNetworkType(e.target.value)}
          value={networkType}
        >
          <Stack direction="row">
            <Radio value="feedforward">Feed Forward</Radio>
            <Radio value="other">Other</Radio>{" "}
            {/* Add other network types as needed */}
          </Stack>
        </RadioGroup>
      </FormControl>
      <FormControl>
        <FormLabel>Number of Entities</FormLabel>
        <NumberInput
          min={10}
          value={spawnCount}
          onChange={(_, value) => setSpawnCount(value)}
        >
          <NumberInputField />
        </NumberInput>
      </FormControl>
      <FormControl>
        <FormLabel>Additional Parameter 1</FormLabel>
        <Input
          placeholder="Enter additional parameter 1"
          value={additionalParam1}
          onChange={(e) => setAdditionalParam1(e.target.value)}
        />
      </FormControl>
      <FormControl>
        <FormLabel>Additional Parameter 2</FormLabel>
        <Input
          placeholder="Enter additional parameter 2"
          value={additionalParam2}
          onChange={(e) => setAdditionalParam2(e.target.value)}
        />
      </FormControl>
      <FormControl>
        <FormLabel>Additional Parameter 3</FormLabel>
        <Input
          placeholder="Enter additional parameter 3"
          value={additionalParam3}
          onChange={(e) => setAdditionalParam3(e.target.value)}
        />
      </FormControl>
      <Button mt={4} colorScheme="blue" type="submit">
        Initialize Training
      </Button>
    </form>
  );
}

export default OTraining;
