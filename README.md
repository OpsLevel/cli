OpsLevel CLI
---

```
opslevel create deploy -s "foo"

cat << EOF | opslevel create deploy -f -
service: "foo"
EOF
```

```
docker run -it --rm public.ecr.aws.com/opslevel/cli:0.0.1 create deploy -s "foo"
```