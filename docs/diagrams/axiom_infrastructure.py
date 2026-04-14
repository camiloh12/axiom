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
from diagrams.aws.integration import SNS, SQS, StepFunctions
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
        ecr      = ECR("ECR\n8 repositories\n(immutable tags)")
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
                gw = ECS("API Gateway\nJWT verify · routing\n256 CPU / 512 MB  ×2–6")

                with Cluster("Application Services"):
                    identity  = ECS("Identity\nauth · RBAC · JWT\n512 CPU / 1 GB  ×2–4")
                    audit     = ECS("Audit Core\nengagements · evidence\n1024 CPU / 2 GB  ×2–8")
                    tb        = ECS("Trial Balance\nTB data · populations\n512 CPU / 1 GB  ×2–4")
                    workpaper = ECS("Workpaper\nYjs real-time sync\n512 CPU / 1 GB  ×2–6")
                    reporting = ECS("Reporting\nasync PDF gen\n512 CPU / 1 GB  ×1–4")
                    doc_proc  = ECS("Doc Processing\nPython + Tesseract\n1024 CPU / 2 GB  ×1–4")

            # Data layer
            with Cluster("Data Layer"):
                rds = RDS(
                    "RDS PostgreSQL 18\ndb.r7g.xlarge  Multi-AZ\n"
                    "identity · core (RLS) · trial_balance\n"
                    "workpaper · reporting\npgvector + pg_stat_statements"
                )
                secrets = SecretsManager(
                    "Secrets Manager\nDB creds · JWT keys\nOAuth secrets\n30-day auto-rotation"
                )
                kms = KMS(
                    "AWS KMS\naxiom-prod-default\naxiom-prod-hipaa\naxiom-prod-rds"
                )

    # ══════════════════════════════════════════════════════════════════════════
    # AWS MANAGED SERVICES  (outside VPC, accessed via VPC endpoints)
    # ══════════════════════════════════════════════════════════════════════════
    with Cluster("AWS Managed Services"):
        bedrock = Sagemaker("Amazon Bedrock\nclaude-haiku-4-5\nclaude-sonnet-4-6\n(via VPC endpoint)")
        sfn     = StepFunctions("Step Functions\nEngagementLifecycle\nDocRequestReminder")
        ses     = SES("SES\naxiom.com\nDKIM + SPF + DMARC")

    with Cluster("SQS Queues  (+ DLQ per queue)"):
        sqs_doc = SQS("document-uploaded\nAudit Core → Audit Core")
        sqs_usr = SQS("user-deactivated\nIdentity → Audit Core")
        sqs_fco = SQS("fco-created\nIdentity → Audit Core")

    with Cluster("S3 Storage"):
        s3_spa      = S3("SPA assets\n(CloudFront OAC\nprivate)")
        s3_evidence = S3("evidence\nSSE-KMS (HIPAA key)\nversioned + cross-region\nreplication (prod)")
        s3_archive  = S3("archive\nObject Lock COMPLIANCE\nWORM  SSE-KMS")
        s3_reports  = S3("reports\nSSE-S3  versioned\nIA after 30d")

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

    # Internet → VPC gateway
    alb >> gw

    # Gateway → services (internal Service Connect)
    gw >> [identity, audit, tb, workpaper, reporting, doc_proc]

    # Services → RDS (PgBouncer sidecar on each service, connects to :5432)
    identity  >> rds
    audit     >> rds
    tb        >> rds
    workpaper >> rds
    reporting >> rds

    # Audit Core → S3
    audit >> s3_evidence
    audit >> s3_archive
    reporting >> s3_reports

    # Async messaging
    audit    >> sqs_doc
    identity >> sqs_usr
    identity >> sqs_fco
    sqs_doc  >> audit
    sqs_usr  >> audit
    sqs_fco  >> audit

    # AI inference (via Bedrock VPC endpoint)
    audit     >> bedrock
    tb        >> bedrock
    workpaper >> bedrock

    # Workflow orchestration
    audit >> sfn

    # Transactional email
    identity >> ses
    audit    >> ses

    # Secrets Manager → services (dashed = startup credential fetch)
    secrets >> Edge(**DASHED) >> gw
    secrets >> Edge(**DASHED) >> identity
    secrets >> Edge(**DASHED) >> audit
    secrets >> Edge(**DASHED) >> tb
    secrets >> Edge(**DASHED) >> workpaper
    secrets >> Edge(**DASHED) >> reporting

    # KMS encryption (dashed = key usage, not data flow)
    kms >> Edge(**DASHED) >> rds
    kms >> Edge(**DASHED) >> s3_evidence
    kms >> Edge(**DASHED) >> s3_archive
    kms >> Edge(**DASHED) >> cw

    # CI/CD pipeline
    github   >> ecr
    github   >> tf_state
    ecr >> Edge(**DASHED) >> gw
    ecr >> Edge(**DASHED) >> identity
    ecr >> Edge(**DASHED) >> audit
    ecr >> Edge(**DASHED) >> tb
    ecr >> Edge(**DASHED) >> workpaper
    ecr >> Edge(**DASHED) >> reporting
    ecr >> Edge(**DASHED) >> doc_proc

    # Telemetry → observability
    gw        >> Edge(**DASHED) >> cw
    identity  >> Edge(**DASHED) >> cw
    audit     >> Edge(**DASHED) >> cw
    tb        >> Edge(**DASHED) >> cw
    workpaper >> Edge(**DASHED) >> cw
    reporting >> Edge(**DASHED) >> cw
    doc_proc  >> Edge(**DASHED) >> cw

    gw        >> Edge(**DASHED) >> xray
    identity  >> Edge(**DASHED) >> xray
    audit     >> Edge(**DASHED) >> xray
    tb        >> Edge(**DASHED) >> xray
    workpaper >> Edge(**DASHED) >> xray
    reporting >> Edge(**DASHED) >> xray

    # Alarms → SNS → ops
    cw >> sns
