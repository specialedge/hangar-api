FROM scratch

ADD hangar-api /bin/

ARG GITSHA="VCS ref not found"
ARG BUILDDATE="Build date not found"
LABEL org.label-schema.vendor="SpecialEdge" \
      org.label-schema.url="https://github.com/specialedge/hangar-api" \
      org.label-schema.name="Hangar-API" \
      org.label-schema.license="GPL v3.0" \
      org.label-schema.vcs-url="https://github.com/specialedge/hangar-api" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.description="Hangar is a distributed artefact management system" \
      org.label-schema.vcs-ref=$GITSHA \
      org.label-schema.build-date=$BUILDDATE

EXPOSE 8080

ENTRYPOINT ["/bin/hangar-api"]