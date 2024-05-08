import os
import io
from PyPDF2 import PdfReader, PdfWriter
from reportlab.pdfgen import canvas
from datetime import datetime
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.serialization import load_pem_private_key, load_pem_public_key
from cryptography.hazmat.backends import default_backend
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import serialization

PRIVATE_KEY_PATH = os.getenv('PRIVATE_KEY_PATH')
PUBLIC_KEY_PATH = os.getenv('PUBLIC_KEY_PATH')
PASSPHRASE = os.getenv('PASSPHRASE').encode() if os.getenv('PASSPHRASE') else None
PDF_PASSWORD = "secret"  # This could also be an environment variable for enhanced security

def add_signature_to_pdf(input_pdf_path, output_pdf_path, image_path, x_position, y_position):
    reader = PdfReader(input_pdf_path)
    writer = PdfWriter()

    for page in reader.pages:
        writer.add_page(page)

    packet = io.BytesIO()
    can = canvas.Canvas(packet, pagesize=reader.pages[0]['/MediaBox'][2:])
    can.drawImage(image_path, x_position, y_position)
    can.save()

    packet.seek(0)
    new_pdf = PdfReader(packet)
    writer.add_page(new_pdf.pages[0])

    with open(output_pdf_path, 'wb') as fout:
        writer.write(fout)

def verify_signature(signature, message, public_key_path):
    with open(public_key_path, "rb") as key_file:
        public_key = load_pem_public_key(
            key_file.read(),
            backend=default_backend()
        )

    try:
        public_key.verify(
            signature,
            message,
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        return True
    except InvalidSignature:
        return False

def sign_message(message, private_key_path, passphrase):
    with open(private_key_path, "rb") as key_file:
        private_key = load_pem_private_key(
            key_file.read(),
            password=passphrase,
            backend=default_backend()
        )

    signature = private_key.sign(
        message,
        padding.PSS(
            mgf=padding.MGF1(hashes.SHA256()),
            salt_length=padding.PSS.MAX_LENGTH
        ),
        hashes.SHA256()
    )
    return signature

def encrypt_pdf(input_pdf_path, output_pdf_path, password):
    reader = PdfReader(input_pdf_path)
    writer = PdfWriter()

    for page in reader.pages:
        writer.add_page(page)

    writer.encrypt(user_pwd=password, owner_pwd=None, use_128bit=True)

    with open(output_pdf_path, 'wb') as fout:
        writer.write(fout)

if __name__ == "__main__":
    input_path = "input.pdf"
    signed_output_path = "signed_output.pdf"
    encrypted_output_path = "encrypted_signed_output.pdf"
    signature_image_path = "signature_image.png"

    add_signature_to_pdf(input_path, signed_output_path, signature_image_path, 100, 100)
    message = b"This is a test message."
    signature = sign_message(message, PRIVATE_KEY_PATH, PASSPHRASE)
    is_verified = verify_signature(signature, message, PUBLIC_KEY_PATH)
    print(f"Signature Verified: {is_verified}")

    # Encrypt the signed PDF
    encrypt_pdf(signed_output_path, encrypted_output_path, PDF_PASSWORD)
    print(f"PDF encrypted successfully. Encrypted file: {encrypted_output_path}")