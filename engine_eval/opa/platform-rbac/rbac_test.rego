package app.rbac

import future.keywords

roles := {
                "tenant_a": {
                    "billing": [
                        {
                            "permission": "read",
                            "resource": "subscriptions"
                        },
                        {
                            "permission": "delete",
                            "resource": "subscriptions"
                        },
                        {
                            "permission": "write",
                            "resource": "subscriptions"
                        }
                    ],
                    "customer": [
                        {
                            "permission": "read",
                            "resource": "settings"
                        },
                        {
                            "permission": "read",
                            "resource": "workspace"
                        },
                        {
                            "permission": "read",
                            "resource": "subscriptions"
                        }
                    ],
                    "employee": [
                        {
                            "permission": "read",
                            "resource": "settings"
                        },
                        {
                            "permission": "read",
                            "resource": "projects"
                        },
                        {
                            "permission": "write",
                            "resource": "settings"
                        },
                        {
                            "permission": "write",
                            "resource": "projects"
                        }
                    ]
                },
                "tenant_b": {
                    "billing": [
                        {
                            "permission": "read",
                            "resource": "subscriptions"
                        },
                        {
                            "permission": "write",
                            "resource": "subscriptions"
                        }
                    ],
                    "customer": [
                        {
                            "permission": "read",
                            "resource": "settings"
                        },
                        {
                            "permission": "read",
                            "resource": "subscriptions"
                        }
                    ],
                    "employee": [
                        {
                            "permission": "read",
                            "resource": "settings"
                        },
                        {
                            "permission": "read",
                            "resource": "projects"
                        },
                        {
                            "permission": "write",
                            "resource": "settings"
                        }
                    ]
                }
            }

user_roles :=  {
                    "alice": [
                        "admin"
                    ],
                    "bob": [
                        "employee",
                        "billing"
                    ],
                    "mark": [
                        "billing"
                    ]
                }


test_allow_delete_with_data_tenant_A {
    allow with input as {
                        "permission": "delete",
                        "resource": "subscriptions",
                        "tenant_id": "tenant_a",
                        "user": "mark"
                    }
    with data.user_roles as user_roles
    with data.roles as roles
}

test_fail_delete_with_data_tenant_B {
    not allow with input as {
                        "permission": "delete",
                        "resource": "subscriptions",
                        "tenant_id": "tenant_b",
                        "user": "mark"
                    }
    with data.user_roles as user_roles
    with data.roles as roles
}