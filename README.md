<div align="center">
  
  ![freebird](https://github.com/sachinsenal0x64/picx-images-hosting/raw/master/freebird.4n7w6lshyp.webp)
 
</div>

# 💕 Community

> 🍻 Join the community:  <a href="https://discord.gg/EbfftZ5Dd4">Discord</a>
> [![](https://cdn.statically.io/gh/sachinsenal0x64/picx-images-hosting@master/discord.72y8nlaw5mdc.webp)](https://discord.gg/EbfftZ5Dd4)
 
<br><br>

> [!NOTE]
> Currently, it's down. I will move my APIs to a new server, which might take some time.

<br><br>

## 🎞️ DEMO

[![hifi tui](https://img.youtube.com/vi/_cwmAqvnJ68/0.jpg)](https://www.youtube.com/watch?v=_cwmAqvnJ68)


# 💨 Quick Start

> Got a torrent file? Convert it to a magnet link: https://nutbread.github.io/t2m/


<br>

## 📡 API DOCUMENTATION (No account required)


------------------------------------------------------------------------------------------

<details>

 <summary><code>GET</code>   <code><b>/v1/magnet</b></code> </summary>

#### Example


<br>

> | Parameter  |   Type    | Description |
> |------------|-----------|-------------|
> | `uri`      |  string   |  magnet = `magnet:?xt` |



<br>


## Request

>xh

    xh GET https://freebird-tumj.onrender.com/v1/magnet?uri="magnet:?xt"

<br>


### Response

  ```json
{
  "chunks": 32,
  "crc": 1,
  "download": "abc",
  "filename": "abc",
  "id": "abc",
  "mimeType": "video/x-matroska",
  "streamable": 1
}

```
<br>


### Status Codes

Freebird returns the following status codes in its API:

> | Status Code | Description |
> | :---        | :--- |
> | 200         | `OK` |
> | 422         | `UNPROCESSABLE CONTENT` |
> | 404         | `NOT FOUND` |
> | 500         | `INTERNAL SERVER ERROR` |


</details>

------------------------------------------------------------------------------------------
