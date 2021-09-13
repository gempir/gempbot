export function getCookie(name: string): string | null {
    return null;

    // const v = document.cookie.match('(^|;) ?' + name + '=([^;]*)(;|$)');
    // return v ? v[2] : null;
}

export function setCookie(name: string, value: string, days: number) {
    const d = new Date();
    d.setTime(d.getTime() + 24 * 60 * 60 * 1000 * days);
    document.cookie = name + "=" + value + ";path=/;expires=" + d.toUTCString();
}

export function deleteCookie(name: string) {
    setCookie(name, '', -1);
}
