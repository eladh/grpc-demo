package org.company.main;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;

@RestController
@EnableAutoConfiguration
public class EchoServer {

    private ObjectMapper mapper = new ObjectMapper();

    @RequestMapping("/")
    String home(@RequestBody JsonJoinToken request) throws IOException {
        JsonAdminToken adminToken = new JsonAdminToken("token-for-" + request.getNodeId() + "-and-" + request.getServiceId());
        return mapper.writeValueAsString(adminToken);
    }

    public static void main(String[] args) {
        SpringApplication.run(EchoServer.class, args);
    }

}