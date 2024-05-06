import React, { Component } from "react";
import axios from "axios";
import {
  Box,
  Flex,
  Wrap,
  WrapItem,
  Icon,
  Text,
  Button,
} from "@chakra-ui/react";
import {
  FiFolder,
  FiSettings,
  FiBookOpen,
  FiLogIn,
  FiLogOut,
  FiEye,
} from "react-icons/fi";

import { FaFolder, FaEye, FaBookOpen } from "react-icons/fa";

import { FaHome } from "react-icons/fa";

import { MdNotStarted } from "react-icons/md";

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link as RLink,
  useParams,
} from "react-router-dom";
import { FaTable } from "react-icons/fa6";
import OFileBrowser from "./OFileBrowser";

import OInitialize from "./OInitialize";
import OModelViewer from "./OModelViewer";
import ODataViewer from "./ODataViewer";

import OMount from "./OMount";

//prod
const protocolPrefix = window.location.protocol === "https:" ? "wss:" : "ws:";
var wsUrl = `${protocolPrefix}//${window.location.host}`;

//testing
wsUrl = "ws://localhost:4123";

//var deployMode = true;

var webSockUrl = wsUrl + "/msg";

class OHome extends Component {
  state = {
    isLoggedIn: false,
    content: "Please log in.",
    dbNameSync: "",
    collectionSync: "",
    tableNameSync: "",
    datasetSize: "",
  };

  componentDidMount() {
    try {
      var ws = new WebSocket(webSockUrl);

      ws.onopen = () => {
        console.log("Connected to server");
      };

      ws.onmessage = (event) => {
        this.wOnmsgProcess(event);
      };

      ws.onclose = (event) => {
        if (event.wasClean) {
          console.log("Disconnected from server cleanly");
        } else {
          console.log("Disconnected from server due to a transmission error");
        }
        console.log("Trying to reconnect to the server...");
        setTimeout(this.connectWebSocket, 3000); // Attempt to reconnect after 3 seconds
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      this.setState({ ws });
    } catch (ex) {
      console.log(ex);
    }
  }

  connectWebSocket = () => {
    var ws = new WebSocket(webSockUrl);

    ws.onopen = () => {
      console.log("Reconnected to server");
    };

    ws.onmessage = (event) => {
      console.log(event);
      //console.log("Message from server: ", event.data);
      this.wOnmsgProcess(event);
    };

    ws.onclose = () => {
      console.log("Disconnected, trying to reconnect...");
      setTimeout(this.connectWebSocket, 3000);
    };

    ws.onerror = (error) => {
      console.error("WebSocket error on reconnect:", error);
    };

    this.setState({ ws });
  };

  wOnmsgProcess(event) {
    try {
      console.log("Message from server: ", event.data);
      var j = JSON.parse(event.data);
      var tmpDBName = "",
        tmpColl = "",
        tmpTblName = "",
        tmpDatasetSize = "";
      var lstk = Object.keys(j);

      if (lstk.includes("dbNameSync")) {
        tmpDBName = j["dbNameSync"];
      }
      if (lstk.includes("collectionSync")) {
        tmpColl = j["collectionSync"];
      }
      if (lstk.includes("tableNameSync")) {
        tmpTblName = j["tableNameSync"];
      }

      if (lstk.includes("datasetSize")) {
        tmpDatasetSize = j["datasetSize"];
      }

      console.log(tmpDBName);
      this.setState({
        dbNameSync: tmpDBName,
        collectionSync: tmpColl,
        tableNameSync: tmpTblName,
        datasetSize: tmpDatasetSize,
      });
    } catch (ex) {
      console.log(ex);
    }
  }

  componentWillUnmount() {
    if (this.state.ws) this.state.ws.close();
  }

  handleDataUpdate = (dbName, collectionName) => {
    // Do something with the received data. You can update state, etc.
    console.log("Received dbName:", dbName);
    console.log("Received collectionName:", collectionName);
    axios
      .post("http://localhost:4123/mountpopulation", {
        dbName,
        collectionName,
      })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.error("Error in POST request:", error);
      });
  };

  handleTrainingDataUpdate = (tableName) => {
    // Do something with the received data. You can update state, etc.
    console.log("Received tableName:", tableName);
    axios
      .post("http://localhost:4123/mounttrainingdata", {
        tableName,
      })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.error("Error in POST request:", error);
      });
  };

  handleLoginLogout = () => {
    this.setState((prevState) => ({
      isLoggedIn: !prevState.isLoggedIn,
      content: prevState.isLoggedIn ? "Please log in." : "Welcome back!",
    }));
  };

  render() {
    const { dbNameSync, collectionSync, tableNameSync, datasetSize } =
      this.state;
    return (
      <Box display="flex" flexDirection="column" h="100vh" w="100vw">
        <Router>
          {/* Top Menu */}
          <Flex
            justifyContent="space-between"
            alignItems="center"
            bg="gray.100"
            p={4}
            wrap="wrap"
            bg="#4A5568"
          >
            {/* Logo on the left */}
            <Text fontSize="xl" fontWeight="bold">
              Cortex Builder
            </Text>

            {/* Center Menu */}
            <Wrap spacing={8} wrap="wrap" justify="center">
              <WrapItem>
                <RLink to="/">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={FaHome} /> Home
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>
              <WrapItem>
                <RLink to="/files">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={FaFolder} /> Files
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>

              <WrapItem>
                <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                  <Text mt={2}>
                    <Icon as={FiSettings} /> Settings
                  </Text>
                </Box>
              </WrapItem>
            </Wrap>

            {/* Login/Logout on the right */}
            <Button
              onClick={this.handleLoginLogout}
              leftIcon={
                this.state.isLoggedIn ? (
                  <Icon as={FiLogOut} />
                ) : (
                  <Icon as={FiLogIn} />
                )
              }
            >
              {this.state.isLoggedIn ? "Logout" : "Login"}
            </Button>
          </Flex>

          <Flex
            justifyContent="space-between"
            alignItems="center"
            bg="gray.100"
            p={4}
            wrap="wrap"
            bg="#6d7d98"
          >
            <Wrap spacing={8} wrap="wrap" justify="center">
              <WrapItem>
                <RLink to="/modelviewer">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={FaEye} /> Model viewer
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>

              <WrapItem>
                <RLink to="/initialize">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={MdNotStarted} /> Initialize
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>

              <WrapItem>
                <RLink to="/dataviewer">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={FaTable} /> Data viewer
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>
            </Wrap>
          </Flex>

          {/* Main content */}
          <Box flex="1" p={4} bg="gray.50">
            {/* Dynamic content based on login status */}

            <Switch>
              <Route exact path="/">
                {this.state.content}
                <OMount
                  dbNameSync={dbNameSync}
                  collectionSync={collectionSync}
                  tableNameSync={tableNameSync}
                  datasetSize={datasetSize}
                />
              </Route>
              <Route exact path="/files">
                <OFileBrowser />
              </Route>

              <Route exact path="/modelviewer">
                <OModelViewer onDataUpdate={this.handleDataUpdate} />
              </Route>

              <Route exact path="/initialize">
                <OInitialize toast={this.props.toast} />
              </Route>

              <Route exact path="/dataviewer">
                <ODataViewer
                  handleTrainingDataUpdate={this.handleTrainingDataUpdate}
                />
              </Route>
            </Switch>
          </Box>
        </Router>
      </Box>
    );
  }
}

// <OTraining  onDataUpdate={this.handleDataUpdate} />

export default OHome;
