#!/usr/bin/env python3
import os
os.environ["PATH"] += os.pathsep + r"C:\Program Files\Graphviz\bin"

"""
Axiom AWS Production Architecture Diagram
Source of truth: docs/specs/infrastructure-design.md

Install deps:
    pip install diagrams
    # Also requires Graphviz: https://graphviz.org/download/
    #   Windows: winget install graphviz  OR  choco install graphviz

Run:
    cd docs/diagrams
    python axiom_infrastructure.py

Output:
    axiom_infrastructure.png  (auto-opens when show=True)
"""

from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import ECS, ECR
from diagrams.aws.database import RDS
from diagrams.aws.engagement import SES
from diagrams.aws.integration import SNS, StepFunctions
from diagrams.aws.management import Cloudtrail, Cloudwatch
from diagrams.aws.ml import Sagemaker  # stand-in for Amazon Bedrock
from diagrams.aws.network import ALB, CloudFront, NATGateway, Route53
from diagrams.aws.security import KMS, SecretsManager, WAF
from diagrams.aws.storage import S3
from diagrams.onprem.vcs import Github

# ── Graph-level Graphviz attributes ──────────────────────────────────────────
GRAPH_ATTR = {
    "fontsize": "13",
    "bgcolor": "#f5f7fa",
    "pad": "1.2",
    "splines": "ortho",
    "ranksep": "1.6",
    "nodesep": "0.9",
    "fontname": "Helvetica Neue",
}

NODE_ATTR = {
    "fontsize": "10",
    "fontname": "Helvetica Neue",
}

# ── Reusable edge styles ──────────────────────────────────────────────────────
DASHED = {"style": "dashed", "color": "#999999"}
SOLID  = {"color": "#444444"}


with Diagram(
    "Axiom — AWS Production Architecture",
    filename="axiom_infrastructure",
    outformat="png",
    graph_attr=GRAPH_ATTR,
    node_attr=NODE_ATTR,
    direction="TB",
    show=False,
):

    # ══════════════════════════════════════════════════════════════════════════
    # INTERNET EDGE — DNS, CDN, WAF
    # ══════════════════════════════════════════════════════════════════════════
    r53 = Route53("Route 53\naxiom.com")

    with Cluster("Edge — WAF + CDN"):
        waf_cf = WAF("WAF\n(CloudFront WebACL)")
        cf     = CloudFront("CloudFront\napp.axiom.com")
        waf_alb = WAF("WAF\n(ALB WebACL)")

    # ══════════════════════════════════════════════════════════════════════════
    # TOOLING ACCOUNT — ECR + CI/CD + Terraform state
    # ══════════════════════════════════════════════════════════════════════════
    with Cluster("axiom-tooling account"):
        github   = Github("GitHub Actions\n(OIDC federation)")
        ecr      = ECR("ECR\n4 repositories\n(immutable tags)")
        tf_state = S3("S3 + DynamoDB\nTerraform state")

    # ══════════════════════════════════════════════════════════════════════════
    # PRODUCTION VPC  us-east-1  10.2.0.0/16
    # ══════════════════════════════════════════════════════════════════════════
    with Cluster("axiom-prod VPC   10.2.0.0/16   us-east-1"):

        # ── Public subnets — ALB + NAT ────────────────────────────────────────
        with Cluster("Public Subnets  (AZ-a  10.2.0.0/25  ·  AZ-b  10.2.0.128/25)"):
            alb   = ALB("ALB\napi.axiom.com :443\nTLS 1.3")
            nat_a = NATGateway("NAT GW\nAZ-a")
            nat_b = NATGateway("NAT GW\nAZ-b")

        # ── Private subnets — ECS + RDS + security ────────────────────────────
        with Cluster("Private Subnets  (AZ-a  10.2.16.0/20  ·  AZ-b  10.2.32.0/20)"):

            # ECS Fargate cluster
            with Cluster("ECS Fargate Cluster  (Service Connect)"):
                axiom_api = ECS(
                    "Axiom API\nGo modular monolith\n"
                    "identity · auditcore · frameworks\n"
                    "workpaper · reporting · ai\n"
                    "1024 CPU / 2 GB  ×2–8"
                )
                doc_proc = ECS(
                    "Doc Processing\nPython + Tesseract\n"
                    "1024 CPU / 2 GB  ×1–4"
                )
                prov_sign = ECS(
                    "Provenance Signer\nGo · KMS Sign-only\n"
                    "512 CPU / 1 GB  ×2–4"
                )

            # Data layer
            with Cluster("Data Layer"):
                rds = RDS(
                    "RDS PostgreSQL 18\ndb.r7g.xlarge  Multi-AZ\n"
                    "axiom_db (RLS all tenants)\n"
                    "pgvector + pg_stat_statements"
                )
                secrets = SecretsManager(
                    "Secrets Manager\nDB creds · JWT keys\nOAuth secrets\n30-day auto-rotation"
                )
                kms = KMS(
                    "AWS KMS\naxiom-prod-default\naxiom-prod-hipaa\naxiom-prod-rds\n"
                    "axiom-prod-provenance-signing\n(ECC_NIST_P256)"
                )

    # ══════════════════════════════════════════════════════════════════════════
    # AWS MANAGED SERVICES  (outside VPC, accessed via VPC endpoints)
    # ══════════════════════════════════════════════════════════════════════════
    with Cluster("AWS Managed Services"):
        bedrock = Sagemaker("Amazon Bedrock\nclaude-haiku-4-5\nclaude-sonnet-4-6\n(via VPC endpoint)")
        sfn     = StepFunctions("Step Functions\nEngagementLifecycle\nDocRequestReminder")
        ses     = SES("SES\naxiom.com\nDKIM + SPF + DMARC")

    with Cluster("S3 Storage"):
        s3_spa      = S3("SPA assets\n(CloudFront OAC\nprivate)")
        s3_evidence = S3("evidence\nSSE-KMS (HIPAA key)\nsigned artifacts\nObject Lock COMPLIANCE")
        s3_archive  = S3("archive\nObject Lock COMPLIANCE\nWORM  SSE-KMS")
        s3_reports  = S3("reports\nObject Lock COMPLIANCE\nSSE-KMS  IA after 30d")
        s3_scf      = S3("scf-catalog\nSCF · OSCAL · AICPA\nCIS crosswalks\nquarterly import")

    # ══════════════════════════════════════════════════════════════════════════
    # OBSERVABILITY
    # ══════════════════════════════════════════════════════════════════════════
    with Cluster("Observability"):
        cw   = Cloudwatch("CloudWatch\nLogs · Metrics\nDashboards · Alarms")
        xray = Cloudwatch("X-Ray\n(OpenTelemetry SDK\n5% sampling prod)")
        ct   = Cloudtrail("CloudTrail\nOrg-wide trail\nS3 + KMS  7yr retention")
        sns  = SNS("SNS\naxiom-prod-ops-alerts\n→ email / PagerDuty")

    # ══════════════════════════════════════════════════════════════════════════
    # CONNECTIONS
    # ══════════════════════════════════════════════════════════════════════════

    # DNS + CDN
    r53     >> waf_cf
    waf_cf  >> cf
    cf      >> s3_spa
    r53     >> waf_alb
    waf_alb >> alb

    # Internet → Axiom API
    alb >> axiom_api

    # Axiom API → Doc Processing (internal HTTP via Service Connect)
    axiom_api >> doc_proc

    # Axiom API → Provenance Signer (internal HTTP via Service Connect)
    axiom_api >> prov_sign

    # Axiom API → RDS (PgBouncer sidecar, connects to :5432)
    axiom_api >> rds

    # Axiom API → S3
    axiom_api >> s3_evidence
    axiom_api >> s3_archive
    axiom_api >> s3_reports
    axiom_api >> s3_scf

    # Provenance Signer → KMS (Sign-only IAM surface)
    prov_sign >> Edge(**DASHED) >> kms
    prov_sign >> s3_evidence

    # AI inference (via Bedrock VPC endpoint)
    axiom_api >> bedrock

    # Workflow orchestration
    axiom_api >> sfn

    # Transactional email
    axiom_api >> ses

    # Secrets Manager → services (dashed = startup credential fetch)
    secrets >> Edge(**DASHED) >> axiom_api

    # KMS encryption (dashed = key usage, not data flow)
    kms >> Edge(**DASHED) >> rds
    kms >> Edge(**DASHED) >> s3_evidence
    kms >> Edge(**DASHED) >> s3_archive
    kms >> Edge(**DASHED) >> s3_reports
    kms >> Edge(**DASHED) >> cw

    # CI/CD pipeline
    github   >> ecr
    github   >> tf_state
    ecr >> Edge(**DASHED) >> axiom_api
    ecr >> Edge(**DASHED) >> doc_proc
    ecr >> Edge(**DASHED) >> prov_sign

    # Telemetry → observability
    axiom_api >> Edge(**DASHED) >> cw
    doc_proc  >> Edge(**DASHED) >> cw
    prov_sign >> Edge(**DASHED) >> cw

    axiom_api >> Edge(**DASHED) >> xray
    prov_sign >> Edge(**DASHED) >> xray

    # Alarms → SNS → ops
    cw >> sns
