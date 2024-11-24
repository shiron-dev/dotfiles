# Ansible

## Tasks

> [!NOTE]
> You can use `xc`(<https://xcfile.dev/>) to run the commands
>
> See <https://xcfile.dev/getting-started/#installation> for installation instructions

### setup

Setup the ansible environment.

Run: once

```bash
ansible-galaxy install -r requirements.yml
```

### lint

Lint the ansible code.

Requires: setup

```bash
ansible-lint
```

### lint:fix

Fix the lint issues in the ansible code.

Requires: setup

```bash
ansible-lint --fix
```

### ansible

Run ansible:check.

Requires: ansible:check

### ansible:check

Dry run the ansible playbook.

Requires: setup

```bash
ansible-playbook -i hosts.yml site.yml -C
```

### ansible:run

Run the ansible playbook.

Requires: setup

```bash
ansible-playbook -i hosts.yml site.yml
```
