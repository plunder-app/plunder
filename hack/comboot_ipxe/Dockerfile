FROM gcc:latest AS IPXE_BUILD
RUN git clone git://git.ipxe.org/ipxe.git
RUN sed -i '/COMBOOT/s/\/\///g' ipxe/src/config/general.h
WORKDIR /ipxe/src/
RUN make bin/undionly.kpxe 

FROM scratch
COPY --from=IPXE_BUILD /ipxe/src/bin/undionly.kpxe .
