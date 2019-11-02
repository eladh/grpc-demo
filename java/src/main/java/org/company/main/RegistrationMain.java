package org.company.main;

import io.grpc.CompressorRegistry;
import io.grpc.DecompressorRegistry;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import org.apache.commons.io.IOUtils;
import org.company.registration.AdminToken;
import org.company.registration.JoinToken;
import org.company.registration.RegistrationGrpc;

import java.io.IOException;
import java.nio.charset.Charset;

/**
 * @author elad hirsch
 */
public class RegistrationMain {
    private final int port;
    private final Server server;


    private static String pubKey;

    private static String prvKey;

    static {
        try {
            pubKey = IOUtils.resourceToString("/pub.cert", Charset.defaultCharset());
            prvKey = IOUtils.resourceToString("/prv.key", Charset.defaultCharset());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private RegistrationMain(ServerBuilder<?> serverBuilder, int port) {
        this.port = port;
        this.server = serverBuilder.addService(new RouteGuideService()).build();
    }

    private Server start() throws IOException, InterruptedException {
        server.start();
        System.out.println("Server started, listening on " + port);
        return server;
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        Server server = new RegistrationMain(ServerBuilder.forPort(2302), 2302).start();

        Server server2 = new RegistrationMain(ServerBuilder.forPort(2303)
                .useTransportSecurity(
                        IOUtils.toInputStream(pubKey, Charset.defaultCharset()),
                        IOUtils.toInputStream(prvKey, Charset.defaultCharset())), 2303).start();
        Server server3 = new RegistrationMain(ServerBuilder.forPort(2304)
                .compressorRegistry(CompressorRegistry.getDefaultInstance())
                .decompressorRegistry(DecompressorRegistry.getDefaultInstance())
                , 2304).start();

        Server server4 = new RegistrationMain(ServerBuilder.forPort(2305)
                .useTransportSecurity(
                        IOUtils.toInputStream(pubKey, Charset.defaultCharset()),
                        IOUtils.toInputStream(prvKey, Charset.defaultCharset()))
                .compressorRegistry(CompressorRegistry.getDefaultInstance())
                .decompressorRegistry(DecompressorRegistry.getDefaultInstance())
                , 2305).start();

        server.awaitTermination();
        server2.awaitTermination();
        server3.awaitTermination();
        server4.awaitTermination();
    }

    private static class RouteGuideService extends RegistrationGrpc.RegistrationImplBase {
        @Override
        public void join(JoinToken request, StreamObserver<AdminToken> responseObserver) {
            responseObserver.onNext(AdminToken.newBuilder().setToken("token-for-" + request.getNodeId() + "-and-" + request.getServiceId()).build());
            responseObserver.onCompleted();
        }
    }
}
