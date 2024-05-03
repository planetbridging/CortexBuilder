import React, { Component } from "react";
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
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
} from "@chakra-ui/react";
import axios from "axios";
import { v4 as uuidv4 } from "uuid";

import OModelViewer from "./OModelViewer";
import ODataViewer from "./ODataViewer";

class OTraining extends Component {
  constructor(props) {
    super(props);
    this.state = {
      networkType: "feedforward",
      spawnCount: 100,
      additionalParam1: "db_" + uuidv4(),
      additionalParam2: "col_" + uuidv4(),
      additionalParam3: "doc_" + uuidv4(),
    };
  }

  onDataChange = (dbName, collectionName) => {
    this.props.onDataUpdate(dbName, collectionName);
  };

  handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const {
        networkType,
        spawnCount,
        additionalParam1,
        additionalParam2,
        additionalParam3,
      } = this.state;
      const response = await axios.post("http://localhost:4123/initialize", {
        networkType,
        spawnCount,
        additionalParam1,
        additionalParam2,
        additionalParam3,
      });
      console.log(response.data);
      this.props.toast({
        title: "Training initialized",
        description:
          "The neural network training initialization process has started.",
        status: "success",
        duration: 5000,
        isClosable: true,
      });
    } catch (error) {
      console.error("Error:", error);
      this.props.toast({
        title: "Error",
        description: "Failed to initialize training.",
        status: "error",
        duration: 5000,
        isClosable: true,
      });
    }
  };

  render() {
    const {
      networkType,
      spawnCount,
      additionalParam1,
      additionalParam2,
      additionalParam3,
    } = this.state;
    return (
      <Tabs>
        <TabList>
          <Tab>Model manager</Tab>
          <Tab>Initialize</Tab>
          <Tab>Evaluation</Tab>
          <Tab>Three</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            <OModelViewer onDataUpdate={this.onDataChange} />
          </TabPanel>
          <TabPanel>
            <form onSubmit={this.handleSubmit}>
              <FormControl>
                <FormLabel>Network Type</FormLabel>
                <RadioGroup
                  onChange={(e) =>
                    this.setState({ networkType: e.target.value })
                  }
                  value={networkType}
                >
                  <Stack direction="row">
                    <Radio value="feedforward">Feed Forward</Radio>
                    <Radio value="other">Other</Radio>
                  </Stack>
                </RadioGroup>
              </FormControl>
              <FormControl>
                <FormLabel>Number of Entities</FormLabel>
                <NumberInput
                  min={10}
                  value={spawnCount}
                  onChange={(_, value) => this.setState({ spawnCount: value })}
                >
                  <NumberInputField />
                </NumberInput>
              </FormControl>
              <FormControl>
                <FormLabel>Database</FormLabel>
                <Input
                  placeholder="Enter additional parameter 1"
                  value={additionalParam1}
                  onChange={(e) =>
                    this.setState({ additionalParam1: e.target.value })
                  }
                />
              </FormControl>
              <FormControl>
                <FormLabel>Additional Parameter 2</FormLabel>
                <Input
                  placeholder="Enter additional parameter 2"
                  value={additionalParam2}
                  onChange={(e) =>
                    this.setState({ additionalParam2: e.target.value })
                  }
                />
              </FormControl>
              <FormControl>
                <FormLabel>Additional Parameter 3</FormLabel>
                <Input
                  placeholder="Enter additional parameter 3"
                  value={additionalParam3}
                  onChange={(e) =>
                    this.setState({ additionalParam3: e.target.value })
                  }
                />
              </FormControl>
              <Button mt={4} colorScheme="blue" type="submit">
                Initialize Training
              </Button>
            </form>
          </TabPanel>
          <TabPanel>
            <ODataViewer />
          </TabPanel>
          <TabPanel>
            <p>three!</p>
          </TabPanel>
        </TabPanels>
      </Tabs>
    );
  }
}

export default OTraining;
