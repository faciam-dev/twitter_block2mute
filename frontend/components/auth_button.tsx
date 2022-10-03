export default function AuthButton({ url }) {
  return (
    <>
      <p>
        <a target="_blank" href="`${url}`" rel="noopener noreferrer">
          認証
        </a>
      </p>
    </>
  );
}
