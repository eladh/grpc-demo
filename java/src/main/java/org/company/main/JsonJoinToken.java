package org.company.main;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * @author elad hirsch
 */
public class JsonJoinToken {
    @JsonProperty("service_id")
    private String serviceId;
    @JsonProperty("node_id")
    private String nodeId;

    public String getNodeId() {
        return nodeId;
    }

    public void setNodeId(String nodeId) {
        this.nodeId = nodeId;
    }

    public String getServiceId() {
        return serviceId;
    }

    public void setServiceId(String serviceId) {
        this.serviceId = serviceId;
    }
}
