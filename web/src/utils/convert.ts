export function formatBigNumber(num: number): string {
  if (num >= 100000000) {
    return (num / 100000000).toFixed(1) + "亿";
  }
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + "万";
  }
  return num.toString();
}

export function formatBytes(bytes: number, decimals: number = 2): string {
  const UNITS = ["B", "KB", "MB", "GB", "TB", "PB"];
  const KILO = 1024;
  if (bytes < 0) {
    return `-${formatBytes(Math.abs(bytes), decimals)}`;
  }
  if (bytes === 0) {
    return "0 B";
  }

  const i = Math.floor(Math.log(bytes) / Math.log(KILO));
  const convertedValue = bytes / Math.pow(KILO, i);

  const fixedDecimals = Math.max(0, decimals);

  return `${convertedValue.toFixed(fixedDecimals)} ${UNITS[i]}`;
}

export function formatDate(timestamp: number): string {
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  const seconds = date.getSeconds().toString().padStart(2, "0");
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
