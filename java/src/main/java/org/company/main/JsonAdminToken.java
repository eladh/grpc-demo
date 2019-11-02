package org.company.main;

/**
 * @author elad hirsch
 */
public class JsonAdminToken {
    private String token;

    public JsonAdminToken(String token) {
        this.token = token;
    }

    public JsonAdminToken() {

    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }
}
