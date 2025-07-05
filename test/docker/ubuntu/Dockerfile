# Dockerfile for testing dotfiles on Ubuntu

FROM ubuntu:latest

# Set DEBIAN_FRONTEND to noninteractive to avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git \
    curl \
    python3 \
    python3-pip \
    sudo \
    make \
    software-properties-common && \
    add-apt-repository --yes --update ppa:ansible/ansible && \
    apt-get install -y --no-install-recommends \
    ansible && \
    rm -rf /var/lib/apt/lists/*

# Install Go (latest stable version)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    golang-go && \
    rm -rf /var/lib/apt/lists/*

# Install gh CLI
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    gh && \
    rm -rf /var/lib/apt/lists/*

# Define the target directory for dotfiles, matching Ansible's default expectation
ENV DOTFILES_DIR /root/projects/github.com/shiron-dev/dotfiles

# Create a working directory and copy the repository
WORKDIR ${DOTFILES_DIR}
# The build context is now the repository root.
COPY . ${DOTFILES_DIR}/

# Install brew-management (Go tool)
# brew-management requires Homebrew, which is not installed in this Docker image by default.
# If brew-management is essential for the setup, Homebrew installation or an alternative approach is needed.
RUN cd scripts/brew-management && go install

# Install Ansible collections
RUN cd scripts/ansible && ansible-galaxy collection install -r requirements.yml

# Run Ansible playbook
# The user running ansible might need sudo privileges for certain tasks.
# We assume the ansible user is root or has passwordless sudo.
# If not, USER directive or sudo configuration might be needed.
# The original setup.bash suggests running brew-management install before ansible.
# RUN brew-management install # This requires Homebrew
# Ansible will use the default home_dir (/root) and dotfiles_dir (from its own facts or playbook variables if not overridden)
# Since we copied the repo to where Ansible expects it by default (based on set_facts logic if user is root),
# no extra variables for paths should be needed here.
RUN cd scripts/ansible && ansible-playbook -i hosts.yml site.yml

# Set a default command (optional, can be overridden)
CMD ["/bin/bash"]
