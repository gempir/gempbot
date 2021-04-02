export function ProfilePicture({ src }: { src: string } = { src: "https://static-cdn.jtvnw.net/user-default-pictures-uv/cdd517fe-def4-11e9-948e-784f43822e80-profile_image-70x70.png" }) {
    return <img src={src} width={"30"} height={"30"} alt="profile" />;
}