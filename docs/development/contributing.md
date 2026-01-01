# Contributing Guidelines

Thank you for your interest in contributing to gSpend!

---

## Getting Started

1. Fork the repository
2. Clone your fork
3. Set up the development environment (see [local-setup.md](../deployment/local-setup.md))
4. Create a feature branch

---

## Branch Naming

| Type | Format | Example |
|------|--------|---------|
| Feature | `feature/<description>` | `feature/add-export-csv` |
| Bug Fix | `bugfix/<description>` | `bugfix/fix-date-picker` |
| Hotfix | `hotfix/<description>` | `hotfix/auth-token-expiry` |

---

## Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance

**Examples:**
```
feat(auth): add user registration endpoint
fix(frontend): correct date picker timezone issue
docs: update API documentation
test(financial): add transaction service tests
```

---

## Pull Request Process

1. **Update tests** - Add or update tests for your changes
2. **Run tests** - Ensure all tests pass (`make test`)
3. **Update documentation** - If applicable
4. **Create PR** - With clear description
5. **Code review** - Address feedback

### PR Title Format
```
<type>(<scope>): <description>
```

---

## Code Review Checklist

- [ ] Code follows project style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No sensitive data in commits
- [ ] Changes are backward compatible

---

## Development Workflow

```
main ◄── develop ◄── feature/your-feature
          │
          ├── bugfix/some-fix
          │
          └── feature/other-feature
```

1. Branch from `develop`
2. Make changes
3. Submit PR to `develop`
4. After review, merge to `develop`
5. `develop` is periodically merged to `main`

---

## Questions?

Open an issue for questions or suggestions.
