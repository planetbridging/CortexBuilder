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
} from "@chakra-ui/react";
import React from "react";

class OMount extends React.Component {
  state = {};

  showEvaluation() {
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
              <ButtonGroup gap="4">
                <Button colorScheme="blackAlpha">Start</Button>
              </ButtonGroup>
            </Stack>
          </CardBody>
        </Card>
      );
    } else {
      return <p>Need to select db & collection & training data set</p>;
    }
  }

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
              <Box>
                <ButtonGroup gap="4">
                  <Button colorScheme="blackAlpha">Evaluate</Button>
                </ButtonGroup>
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

              <Box>
                <ButtonGroup gap="4">
                  <Button colorScheme="blackAlpha">Placeholder</Button>
                </ButtonGroup>
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
