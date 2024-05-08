import {
  Box,
  Button,
  ButtonGroup,
  Card,
  CardBody,
  CardHeader,
  Heading,
  Stack,
  StackDivider,
  Text,
  Wrap,
  WrapItem,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
} from "@chakra-ui/react";
import React from "react";
import axios from "axios";

class OMount extends React.Component {
  state = {
    batch: 10,
    datasetSize: 1,
  };

  showEvaluation() {
    const { batch, datasetSize } = this.state;
    var dbNameSync = this.props.dbNameSync;
    var collectionSync = this.props.collectionSync;
    var tableNameSync = this.props.tableNameSync;
    if (
      !dbNameSync.startsWith("old_") &&
      !dbNameSync.startsWith("old_") &&
      !collectionSync.startsWith("old_") &&
      !collectionSync.startsWith("old_") &&
      !tableNameSync.startsWith("old_") &&
      !tableNameSync.startsWith("new")
    ) {
      return (
        <Card>
          <CardHeader>
            <Heading size="md">Evaluation</Heading>
          </CardHeader>

          <CardBody>
            <Stack divider={<StackDivider />} spacing="4">
              <Text pt="2" fontSize="sm">
                Batch sizes - {batch}/{this.props.modelCount}
              </Text>
              <NumberInput
                value={batch}
                onChange={(valueString) => this.handleAmountChange(valueString)}
                max={this.props.modelCount}
                min={1}
                clampValueOnBlur={false}
              >
                <NumberInputField placeholder="Amount" />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
              <Text pt="2" fontSize="sm">
                Data sizes - {datasetSize}/{this.props.datasetSize}
              </Text>
              <NumberInput
                value={datasetSize}
                onChange={(valueString) =>
                  this.handleAmountDatasetChange(valueString)
                }
                max={this.props.datasetSize}
                min={1}
                clampValueOnBlur={false}
              >
                <NumberInputField placeholder="Amount" />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
              <ButtonGroup gap="4">
                <Button
                  colorScheme="blackAlpha"
                  onClick={() => this.startEvaluation(batch)}
                >
                  Start(multi unstable)
                </Button>
                <Button
                  colorScheme="blackAlpha"
                  onClick={() => this.startEvaluation(1)}
                >
                  Start single
                </Button>
              </ButtonGroup>
            </Stack>
          </CardBody>
        </Card>
      );
    } else {
      return <p>Need to select db & collection & training data set</p>;
    }
  }

  startEvaluation = async (batch) => {
    const { dbNameSync, collectionSync, tableNameSync, modelCount } =
      this.props;

    var batchNumber = batch;
    batchNumber = batchNumber.toString();
    try {
      const response = await axios.post(
        "http://" + this.props.currentHost + ":4123/evaluation",
        {
          dbName: dbNameSync,
          collectionName: collectionSync,
          tableName: tableNameSync,
          batchNumber,
        }
      );

      // Handle the response (e.g., update UI, show success message, etc.)
      console.log("Evaluation response:", response.data);
    } catch (error) {
      // Handle errors (e.g., show error message)
      console.error("Error during evaluation:", error);
    }
  };

  handleAmountChange = (valueString) => {
    // Convert the string value to a number
    var value = parseInt(valueString, 10) || 0;
    if (value >= this.props.modelCount) {
      value = this.props.modelCount;
    }
    this.setState({ batch: value });
  };

  handleAmountDatasetChange = (valueString) => {
    // Convert the string value to a number
    var value = parseInt(valueString, 10) || 0;
    if (value >= this.props.datasetSize) {
      value = this.props.datasetSize;
    }
    this.setState({ datasetSize: value });
  };

  showSelectedPopulation() {
    var dbNameSync = this.props.dbNameSync;
    var collectionSync = this.props.collectionSync;

    if (
      !dbNameSync.startsWith("old_") &&
      !dbNameSync.startsWith("old_") &&
      !collectionSync.startsWith("old_") &&
      !collectionSync.startsWith("old_")
    ) {
      return (
        <Card>
          <CardHeader>
            <Heading size="md">Selected population</Heading>
          </CardHeader>

          <CardBody>
            <Stack divider={<StackDivider />} spacing="4">
              <Box>
                <Heading size="xs" textTransform="uppercase">
                  DB
                </Heading>
                <Text pt="2" fontSize="sm">
                  {dbNameSync}
                </Text>
              </Box>
              <Box>
                <Heading size="xs" textTransform="uppercase">
                  Collection
                </Heading>
                <Text pt="2" fontSize="sm">
                  {collectionSync}
                </Text>
              </Box>
            </Stack>
          </CardBody>
        </Card>
      );
    } else {
      return <p>No population selected</p>;
    }
  }

  showSelectedDataset() {
    var tableNameSync = this.props.tableNameSync;
    if (!tableNameSync.startsWith("old_") && !tableNameSync.startsWith("new")) {
      return (
        <Card>
          <CardHeader>
            <Heading size="md">Selected Dataset</Heading>
          </CardHeader>

          <CardBody>
            <Stack divider={<StackDivider />} spacing="4">
              <Box>
                <Heading size="xs" textTransform="uppercase">
                  Dataset
                </Heading>
                <Text pt="2" fontSize="sm">
                  {tableNameSync}
                </Text>
              </Box>
            </Stack>
          </CardBody>
        </Card>
      );
    } else {
      return <p>No selected dataset</p>;
    }
  }

  render() {
    var selectedPop = this.showSelectedPopulation();
    var selectedDataset = this.showSelectedDataset();
    var showEvaluation = this.showEvaluation();

    return (
      <Wrap>
        <WrapItem>{selectedPop}</WrapItem>
        <WrapItem>{selectedDataset}</WrapItem>
        <WrapItem>{showEvaluation}</WrapItem>
      </Wrap>
    );
  }
}

export default OMount;
