import React, { Component } from "react";
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

class OHome extends Component {
  state = {
    isLoggedIn: false,
    content: "Please log in.",
  };

  handleDataUpdate = (dbName, collectionName) => {
    // Do something with the received data. You can update state, etc.
    console.log("Received dbName:", dbName);
    console.log("Received collectionName:", collectionName);
  };

  handleLoginLogout = () => {
    this.setState((prevState) => ({
      isLoggedIn: !prevState.isLoggedIn,
      content: prevState.isLoggedIn ? "Please log in." : "Welcome back!",
    }));
  };

  render() {
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
                <ODataViewer />
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
