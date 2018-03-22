FROM debian
COPY ./pmd-crd /pmd-crd
ENTRYPOINT /pmd-crd
