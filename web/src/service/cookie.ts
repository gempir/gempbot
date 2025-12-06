export function getCookie(name: string): string | null {
  const v = document.cookie.match(`(^|;) ?${name}=([^;]*)(;|$)`);
  return v ? v[2] : null;
}

export function setCookie(name: string, value: string, days: number = 365) {
  const d = new Date();
  d.setTime(d.getTime() + 24 * 60 * 60 * 1000 * days);
  document.cookie = `${name}=${value};path=/;expires=${d.toUTCString()}`;
}

export function deleteCookie(name: string) {
  setCookie(name, "", -1);
}

export function parseCookie(str: string): Record<string, string> {
  if (str.trim() === "") {
    return {};
  }

  return str
    .split(";")
    .map((v) => v.split("="))
    .reduce((acc: any, v) => {
      acc[decodeURIComponent(v[0].trim())] = decodeURIComponent(v[1].trim());
      return acc;
    }, {});
}
