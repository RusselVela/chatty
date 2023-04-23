ARG RUN_FROM=gcr.io/distroless/static-debian11@sha256:21d3f84a4f37c36199fd07ad5544dcafecc17776e3f3628baf9a57c8c0181b3f
FROM ${RUN_FROM} as image
USER nonroot:nonroot

COPY output/bin/chatty /bin

ENTRYPOINT [ "/bin/chatty" ]