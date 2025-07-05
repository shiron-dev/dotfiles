# Makefile for building and testing dotfiles with Docker

.PHONY: build-docker-ubuntu test-docker-ubuntu test-docker-all test-all

# Makefile for building and testing dotfiles with Docker

.PHONY: build-docker-ubuntu test-docker-ubuntu test-docker-all test-all

# Build the Docker image for Ubuntu
# The build context is set to the repository root (the final '.')
# The Dockerfile is specified with -f
build-docker-ubuntu:
	docker build -t ubuntu-dotfiles-test -f ./test/docker/ubuntu/Dockerfile .

# Test with the Ubuntu Docker image
# Checks if some key symlinks are created by Ansible.
test-docker-ubuntu: build-docker-ubuntu
	docker run --rm ubuntu-dotfiles-test bash -c "\
		ls -l /root/.zshrc && \
		ls -l /root/.gitconfig && \
		ls -l /root/.config/gh/config.yml && \
		echo 'Docker tests passed!'"

# Test with all Docker images (currently only Ubuntu)
test-docker-all: test-docker-ubuntu

# Run all tests (currently only Docker tests)
test-all: test-docker-all
