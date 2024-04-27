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

import { FaFolder, FaEye } from "react-icons/fa";

import { FaHome } from "react-icons/fa";

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link as RLink,
  useParams,
} from "react-router-dom";
import OFileBrowser from "./OFileBrowser";
import OModelViewer from "./OModelViewer";

class OHome extends Component {
  state = {
    isLoggedIn: false,
    content: "Please log in.",
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
                <RLink to="/modelviewer">
                  <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                    <Text mt={2}>
                      <Icon as={FaEye} /> Model Viewer
                    </Text>
                  </Box>
                </RLink>
              </WrapItem>
              <WrapItem>
                <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                  <Text mt={2}>
                    <Icon as={FiBookOpen} /> Training
                  </Text>
                </Box>
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
                <OModelViewer />
              </Route>
            </Switch>
          </Box>
        </Router>
      </Box>
    );
  }
}

export default OHome;
